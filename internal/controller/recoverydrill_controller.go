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

	corev1 "k8s.io/api/core/v1"
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
	"github.com/Kitio-Tek/hilios-operator/internal/finalizers"
	"github.com/Kitio-Tek/hilios-operator/internal/labels"
	"github.com/Kitio-Tek/hilios-operator/internal/metrics"
	"github.com/Kitio-Tek/hilios-operator/internal/predicates"
)

// RecoveryDrillReconciler reconciles a RecoveryDrill object.
type RecoveryDrillReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=resilience.hilios.io,resources=recoverydrills,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=resilience.hilios.io,resources=recoverydrills/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=resilience.hilios.io,resources=recoverydrills/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;create;delete
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile drives a RecoveryDrill through Pending -> Running -> Succeeded/Failed.
// External verification engines (such as Velero restores) plug in by replacing
// the executeDrill body without restructuring the loop.
func (r *RecoveryDrillReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	drill := &resiliencev1alpha1.RecoveryDrill{}
	if err := r.Get(ctx, req.NamespacedName, drill); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("get drill: %w", err)
	}

	if !drill.DeletionTimestamp.IsZero() {
		return r.handleDelete(ctx, drill)
	}

	if finalizers.Add(&drill.ObjectMeta, finalizers.HiliosFinalizer) {
		if err := r.Update(ctx, drill); err != nil {
			return ctrl.Result{}, fmt.Errorf("add finalizer: %w", err)
		}
		return ctrl.Result{Requeue: true}, nil
	}

	switch drill.Status.Phase {
	case "":
		return r.transitionToPending(ctx, drill)
	case resiliencev1alpha1.DrillPhasePending:
		return r.transitionToRunning(ctx, drill)
	case resiliencev1alpha1.DrillPhaseRunning:
		return r.executeDrill(ctx, drill)
	case resiliencev1alpha1.DrillPhaseSucceeded, resiliencev1alpha1.DrillPhaseFailed:
		logger.V(1).Info("drill is terminal", "phase", drill.Status.Phase)
		return ctrl.Result{}, nil
	default:
		return ctrl.Result{}, fmt.Errorf("unknown phase %q", drill.Status.Phase)
	}
}

func (r *RecoveryDrillReconciler) transitionToPending(ctx context.Context, drill *resiliencev1alpha1.RecoveryDrill) (ctrl.Result, error) {
	drill.Status.Phase = resiliencev1alpha1.DrillPhasePending
	drill.Status.ObservedGeneration = drill.Generation
	conditions.True(&drill.Status.Conditions, resiliencev1alpha1.ConditionScheduled,
		resiliencev1alpha1.ReasonScheduled, "drill scheduled", drill.Generation)
	if err := r.Status().Update(ctx, drill); err != nil {
		return ctrl.Result{}, fmt.Errorf("status update: %w", err)
	}
	events.Normal(r.Recorder, drill, resiliencev1alpha1.ReasonScheduled,
		"recovery drill scheduled (type=%s)", drill.Spec.Type)
	return ctrl.Result{Requeue: true}, nil
}

func (r *RecoveryDrillReconciler) transitionToRunning(ctx context.Context, drill *resiliencev1alpha1.RecoveryDrill) (ctrl.Result, error) {
	now := metav1.NewTime(time.Now())
	drill.Status.Phase = resiliencev1alpha1.DrillPhaseRunning
	drill.Status.StartTime = &now
	conditions.True(&drill.Status.Conditions, resiliencev1alpha1.ConditionRunning,
		resiliencev1alpha1.ReasonStarted, "drill running", drill.Generation)
	drill.Status.Evidence = append(drill.Status.Evidence, resiliencev1alpha1.EvidenceRecord{
		Step:    "start",
		Time:    now,
		Result:  "Pass",
		Message: fmt.Sprintf("started %s drill", drill.Spec.Type),
	})
	if err := r.Status().Update(ctx, drill); err != nil {
		return ctrl.Result{}, fmt.Errorf("status update: %w", err)
	}
	events.Normal(r.Recorder, drill, resiliencev1alpha1.ReasonStarted, "drill started")
	return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
}

func (r *RecoveryDrillReconciler) executeDrill(ctx context.Context, drill *resiliencev1alpha1.RecoveryDrill) (ctrl.Result, error) {
	if drill.Status.StartTime != nil {
		elapsed := time.Since(drill.Status.StartTime.Time)
		if elapsed > time.Duration(drill.Spec.TimeoutSeconds)*time.Second {
			return r.complete(ctx, drill, false, resiliencev1alpha1.ReasonTimeoutExceeded,
				fmt.Sprintf("drill exceeded timeout of %ds", drill.Spec.TimeoutSeconds))
		}
	}

	if drill.Spec.Type == resiliencev1alpha1.DrillRestoreVerification {
		if err := r.ensureVerificationNamespace(ctx, drill); err != nil {
			return r.complete(ctx, drill, false, resiliencev1alpha1.ReasonRestoreFailed, err.Error())
		}
	}

	return r.complete(ctx, drill, true, resiliencev1alpha1.ReasonRestoreVerified,
		fmt.Sprintf("drill %s completed", drill.Spec.Type))
}

func (r *RecoveryDrillReconciler) ensureVerificationNamespace(ctx context.Context, drill *resiliencev1alpha1.RecoveryDrill) error {
	name := drill.Spec.VerificationNamespace
	if name == "" {
		name = fmt.Sprintf("hilios-verify-%s", drill.Name)
	}
	ns := &corev1.Namespace{}
	if err := r.Get(ctx, client.ObjectKey{Name: name}, ns); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		ns = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Labels:      labels.MergeManagedBy(map[string]string{labels.LabelDrill: drill.Name}),
				Annotations: map[string]string{labels.AnnotationVerificationNamespace: drill.Name},
			},
		}
		if err := r.Create(ctx, ns); err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("create verification namespace: %w", err)
		}
	}
	return nil
}

func (r *RecoveryDrillReconciler) complete(ctx context.Context, drill *resiliencev1alpha1.RecoveryDrill, success bool, reason, msg string) (ctrl.Result, error) {
	now := metav1.NewTime(time.Now())
	drill.Status.CompletionTime = &now
	if drill.Status.StartTime != nil {
		drill.Status.DurationSeconds = int32(now.Sub(drill.Status.StartTime.Time).Seconds())
	}
	drill.Status.Message = msg

	if success {
		drill.Status.Phase = resiliencev1alpha1.DrillPhaseSucceeded
		conditions.False(&drill.Status.Conditions, resiliencev1alpha1.ConditionRunning,
			resiliencev1alpha1.ReasonCompleted, "drill completed", drill.Generation)
		conditions.True(&drill.Status.Conditions, resiliencev1alpha1.ConditionSucceeded,
			reason, msg, drill.Generation)
		drill.Status.Evidence = append(drill.Status.Evidence, resiliencev1alpha1.EvidenceRecord{
			Step: "complete", Time: now, Result: "Pass", Message: msg,
		})
		events.Normal(r.Recorder, drill, reason, msg)
	} else {
		drill.Status.Phase = resiliencev1alpha1.DrillPhaseFailed
		conditions.False(&drill.Status.Conditions, resiliencev1alpha1.ConditionRunning,
			reason, msg, drill.Generation)
		conditions.True(&drill.Status.Conditions, resiliencev1alpha1.ConditionFailed,
			reason, msg, drill.Generation)
		drill.Status.Evidence = append(drill.Status.Evidence, resiliencev1alpha1.EvidenceRecord{
			Step: "complete", Time: now, Result: "Fail", Message: msg,
		})
		events.Warning(r.Recorder, drill, reason, msg)
	}

	if err := r.Status().Update(ctx, drill); err != nil {
		return ctrl.Result{}, fmt.Errorf("status update: %w", err)
	}

	result := "Failed"
	if success {
		result = "Succeeded"
	}
	metrics.DrillsTotal.WithLabelValues(string(drill.Spec.Type), result).Inc()
	if drill.Status.DurationSeconds > 0 {
		metrics.DrillDurationSeconds.WithLabelValues(string(drill.Spec.Type), result).Observe(float64(drill.Status.DurationSeconds))
	}
	return ctrl.Result{}, nil
}

func (r *RecoveryDrillReconciler) handleDelete(ctx context.Context, drill *resiliencev1alpha1.RecoveryDrill) (ctrl.Result, error) {
	if !finalizers.Has(&drill.ObjectMeta, finalizers.HiliosFinalizer) {
		return ctrl.Result{}, nil
	}
	if drill.Spec.Cleanup && drill.Spec.VerificationNamespace != "" {
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: drill.Spec.VerificationNamespace}}
		if err := r.Delete(ctx, ns); err != nil && !apierrors.IsNotFound(err) {
			return ctrl.Result{}, fmt.Errorf("delete verification namespace: %w", err)
		}
	}
	finalizers.Remove(&drill.ObjectMeta, finalizers.HiliosFinalizer)
	if err := r.Update(ctx, drill); err != nil {
		return ctrl.Result{}, fmt.Errorf("remove finalizer: %w", err)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager registers the reconciler with the supplied manager.
func (r *RecoveryDrillReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor("hilios-recoverydrill")
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&resiliencev1alpha1.RecoveryDrill{}, builder.WithPredicates(predicates.GenerationOrPause())).
		Complete(r)
}
