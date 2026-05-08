/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
	"github.com/Kitio-Tek/hilios-operator/internal/conditions"
	"github.com/Kitio-Tek/hilios-operator/internal/events"
	"github.com/Kitio-Tek/hilios-operator/internal/labels"
	"github.com/Kitio-Tek/hilios-operator/internal/metrics"
	"github.com/Kitio-Tek/hilios-operator/internal/policy"
	"github.com/Kitio-Tek/hilios-operator/internal/predicates"
	"github.com/Kitio-Tek/hilios-operator/internal/scheduling"
)

// ResiliencePolicyReconciler reconciles a ResiliencePolicy object.
type ResiliencePolicyReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=resilience.hilios.io,resources=resiliencepolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=resilience.hilios.io,resources=resiliencepolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=resilience.hilios.io,resources=resiliencepolicies/finalizers,verbs=update
// +kubebuilder:rbac:groups=resilience.hilios.io,resources=recoverydrills,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile evaluates a ResiliencePolicy: it counts matched workloads, validates
// the policy spec, updates status conditions, and (when authorised) creates
// RecoveryDrill objects for the verifications declared in spec.
func (r *ResiliencePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	policy := &resiliencev1alpha1.ResiliencePolicy{}
	if err := r.Get(ctx, req.NamespacedName, policy); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("get policy: %w", err)
	}

	if policy.Spec.Suspend || policy.Annotations[labels.AnnotationPaused] == "true" {
		conditions.False(&policy.Status.Conditions, resiliencev1alpha1.ConditionReady,
			resiliencev1alpha1.ReasonSuspended, "policy is suspended", policy.Generation)
		policy.Status.ObservedGeneration = policy.Generation
		return ctrl.Result{}, r.updateStatus(ctx, policy)
	}

	if err := r.validate(policy); err != nil {
		conditions.False(&policy.Status.Conditions, resiliencev1alpha1.ConditionValidated,
			resiliencev1alpha1.ReasonValidationFailed, err.Error(), policy.Generation)
		conditions.False(&policy.Status.Conditions, resiliencev1alpha1.ConditionReady,
			resiliencev1alpha1.ReasonValidationFailed, err.Error(), policy.Generation)
		policy.Status.ObservedGeneration = policy.Generation
		events.Warning(r.Recorder, policy, resiliencev1alpha1.ReasonValidationFailed, err.Error())
		return ctrl.Result{}, r.updateStatus(ctx, policy)
	}
	conditions.True(&policy.Status.Conditions, resiliencev1alpha1.ConditionValidated,
		resiliencev1alpha1.ReasonReady, "policy spec validated", policy.Generation)

	matched, err := r.countMatchedTargets(ctx, policy)
	if err != nil {
		logger.Error(err, "list matched workloads")
		return ctrl.Result{}, err
	}
	policy.Status.MatchedTargets = matched

	staleCount, err := r.countStaleVerifications(ctx, policy, time.Now())
	if err != nil {
		logger.Error(err, "compute stale verifications")
		return ctrl.Result{}, err
	}
	policy.Status.LastDriftCount = staleCount
	if staleCount > 0 {
		policy.Status.LastViolation = fmt.Sprintf("%d verification(s) past freshness window", staleCount)
	} else {
		policy.Status.LastViolation = ""
	}

	now := metav1.NewTime(time.Now())
	policy.Status.LastEvaluationTime = &now
	policy.Status.ObservedGeneration = policy.Generation

	conditions.True(&policy.Status.Conditions, resiliencev1alpha1.ConditionReady,
		resiliencev1alpha1.ReasonReady, fmt.Sprintf("%d workloads governed", matched), policy.Generation)

	if matched == 0 {
		conditions.True(&policy.Status.Conditions, resiliencev1alpha1.ConditionDegraded,
			resiliencev1alpha1.ReasonSelectorEmpty,
			"selector matched no workloads", policy.Generation)
	} else {
		conditions.False(&policy.Status.Conditions, resiliencev1alpha1.ConditionDegraded,
			resiliencev1alpha1.ReasonReady, "policy targets are governed", policy.Generation)
	}

	metrics.PolicyEvaluationsTotal.WithLabelValues(policy.Namespace, policy.Name).Inc()
	metrics.PolicyDriftCount.WithLabelValues(policy.Namespace, policy.Name).Set(float64(policy.Status.LastDriftCount))

	if err := r.updateStatus(ctx, policy); err != nil {
		return ctrl.Result{}, err
	}

	requeue := scheduling.NextRequeue(policy.Spec.Schedule, time.Now(), 5*time.Minute)
	return ctrl.Result{RequeueAfter: requeue}, nil
}

// validate enforces invariants that cannot be expressed via CRD validation alone.
func (r *ResiliencePolicyReconciler) validate(p *resiliencev1alpha1.ResiliencePolicy) error {
	if len(p.Spec.Verifications) == 0 {
		return fmt.Errorf("at least one verification must be declared")
	}
	if p.Spec.SLO.RecoveryTimeSeconds < 0 {
		return fmt.Errorf("slo.recoveryTimeSeconds must be non-negative")
	}
	for i, v := range p.Spec.Verifications {
		if v.Kind == "" {
			return fmt.Errorf("verifications[%d].kind is required", i)
		}
	}
	return nil
}

func (r *ResiliencePolicyReconciler) countMatchedTargets(ctx context.Context, p *resiliencev1alpha1.ResiliencePolicy) (int32, error) {
	selector, err := metav1.LabelSelectorAsSelector(&p.Spec.TargetSelector)
	if err != nil {
		return 0, err
	}
	if selector.Empty() {
		return 0, nil
	}
	list := &appsv1.StatefulSetList{}
	if err := r.List(ctx, list, &client.ListOptions{LabelSelector: selector}); err != nil {
		return 0, err
	}
	return int32(len(list.Items)), nil
}

// countStaleVerifications joins the policy's verifications with the timestamps
// of the most recent successful RecoveryDrill of the same kind in the same
// namespace, and counts how many have aged past the freshness window.
func (r *ResiliencePolicyReconciler) countStaleVerifications(ctx context.Context, p *resiliencev1alpha1.ResiliencePolicy, now time.Time) (int32, error) {
	drills := &resiliencev1alpha1.RecoveryDrillList{}
	if err := r.List(ctx, drills, client.InNamespace(p.Namespace)); err != nil {
		return 0, err
	}
	lastByType := make(map[resiliencev1alpha1.VerificationKind]time.Time)
	for _, d := range drills.Items {
		if d.Status.Phase != resiliencev1alpha1.DrillPhaseSucceeded || d.Status.CompletionTime == nil {
			continue
		}
		var kind resiliencev1alpha1.VerificationKind
		switch d.Spec.Type {
		case resiliencev1alpha1.DrillRestoreVerification:
			kind = resiliencev1alpha1.VerificationRestoreVerification
		case resiliencev1alpha1.DrillFailoverDrill:
			kind = resiliencev1alpha1.VerificationFailoverDrill
		default:
			continue
		}
		t := d.Status.CompletionTime.Time
		if cur, ok := lastByType[kind]; !ok || t.After(cur) {
			lastByType[kind] = t
		}
	}
	statuses := make([]policy.VerificationStatus, 0, len(p.Spec.Verifications))
	for _, v := range p.Spec.Verifications {
		statuses = append(statuses, policy.VerificationStatus{
			Spec:        v,
			LastSuccess: lastByType[v.Kind],
		})
	}
	return policy.CountStale(statuses, now), nil
}

func (r *ResiliencePolicyReconciler) updateStatus(ctx context.Context, p *resiliencev1alpha1.ResiliencePolicy) error {
	if err := r.Status().Update(ctx, p); err != nil {
		return fmt.Errorf("update policy status: %w", err)
	}
	return nil
}

// SetupWithManager registers the reconciler with the supplied manager.
func (r *ResiliencePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor("hilios-resiliencepolicy")
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&resiliencev1alpha1.ResiliencePolicy{}, builder.WithPredicates(predicates.GenerationOrPause())).
		Owns(&resiliencev1alpha1.RecoveryDrill{}).
		Complete(r)
}
