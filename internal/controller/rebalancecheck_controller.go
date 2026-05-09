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
	"github.com/Kitio-Tek/hilios-operator/internal/metrics"
	"github.com/Kitio-Tek/hilios-operator/internal/predicates"
	"github.com/Kitio-Tek/hilios-operator/internal/safeint"
	"github.com/Kitio-Tek/hilios-operator/internal/scheduling"
	"github.com/Kitio-Tek/hilios-operator/internal/selector"
	"github.com/Kitio-Tek/hilios-operator/internal/topology"
)

// RebalanceCheckReconciler reconciles a RebalanceCheck object.
type RebalanceCheckReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=resilience.hilios.io,resources=rebalancechecks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=resilience.hilios.io,resources=rebalancechecks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=resilience.hilios.io,resources=rebalancechecks/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch

// Reconcile evaluates the topology distribution of pods matching the selector
// and reports drift through standard conditions.
func (r *RebalanceCheckReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	check := &resiliencev1alpha1.RebalanceCheck{}
	if err := r.Get(ctx, req.NamespacedName, check); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("get check: %w", err)
	}

	sel, err := selector.Build(check.Spec.TargetSelector)
	if err != nil {
		conditions.False(&check.Status.Conditions, resiliencev1alpha1.ConditionBalanced,
			resiliencev1alpha1.ReasonValidationFailed, err.Error(), check.Generation)
		return ctrl.Result{}, r.Status().Update(ctx, check)
	}

	pods := &corev1.PodList{}
	if err := r.List(ctx, pods, &client.ListOptions{LabelSelector: sel}); err != nil {
		return ctrl.Result{}, fmt.Errorf("list pods: %w", err)
	}

	nodes := &corev1.NodeList{}
	if err := r.List(ctx, nodes); err != nil {
		return ctrl.Result{}, fmt.Errorf("list nodes: %w", err)
	}

	dist := topology.Distribute(pods.Items, nodes.Items, check.Spec.TopologyKey)
	skew := topology.Skew(dist)
	logger.V(1).Info("distribution", "skew", skew, "domains", len(dist))

	check.Status.MatchedTargets = safeint.Int32(len(pods.Items))
	check.Status.LastSkew = skew
	check.Status.ObservedGeneration = check.Generation
	now := metav1.NewTime(time.Now())
	check.Status.LastEvaluationTime = &now
	check.Status.Distribution = check.Status.Distribution[:0]
	for _, k := range topology.SortedDomains(dist) {
		check.Status.Distribution = append(check.Status.Distribution,
			resiliencev1alpha1.TopologyDistribution{Domain: k, Replicas: dist[k]})
	}

	if skew > check.Spec.MaxSkew {
		conditions.False(&check.Status.Conditions, resiliencev1alpha1.ConditionBalanced,
			resiliencev1alpha1.ReasonReplicaSkew,
			fmt.Sprintf("skew %d exceeds maxSkew %d", skew, check.Spec.MaxSkew), check.Generation)
		conditions.True(&check.Status.Conditions, resiliencev1alpha1.ConditionDrifted,
			resiliencev1alpha1.ReasonReplicaSkew,
			fmt.Sprintf("skew %d exceeds maxSkew %d", skew, check.Spec.MaxSkew), check.Generation)
		conditions.True(&check.Status.Conditions, resiliencev1alpha1.ConditionActionRequired,
			resiliencev1alpha1.ReasonReplicaSkew,
			"apply topology spread or move replicas", check.Generation)
		check.Status.Message = fmt.Sprintf("replica skew detected: %d (max %d)", skew, check.Spec.MaxSkew)
		events.Warning(r.Recorder, check, resiliencev1alpha1.ReasonReplicaSkew, check.Status.Message)
	} else {
		conditions.True(&check.Status.Conditions, resiliencev1alpha1.ConditionBalanced,
			resiliencev1alpha1.ReasonReplicaBalanced,
			fmt.Sprintf("skew %d within max %d", skew, check.Spec.MaxSkew), check.Generation)
		conditions.False(&check.Status.Conditions, resiliencev1alpha1.ConditionDrifted,
			resiliencev1alpha1.ReasonReplicaBalanced, "no drift detected", check.Generation)
		conditions.False(&check.Status.Conditions, resiliencev1alpha1.ConditionActionRequired,
			resiliencev1alpha1.ReasonReplicaBalanced, "no action required", check.Generation)
		check.Status.Message = "topology balanced"
	}

	metrics.RebalanceSkew.WithLabelValues(check.Namespace, check.Name).Set(float64(skew))

	if err := r.Status().Update(ctx, check); err != nil {
		return ctrl.Result{}, fmt.Errorf("status update: %w", err)
	}
	requeue := scheduling.NextRequeue(check.Spec.Schedule, time.Now(), 2*time.Minute)
	return ctrl.Result{RequeueAfter: requeue}, nil
}

// SetupWithManager registers the reconciler with the supplied manager.
func (r *RebalanceCheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor("hilios-rebalancecheck")
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&resiliencev1alpha1.RebalanceCheck{}, builder.WithPredicates(predicates.GenerationOrPause())).
		Complete(r)
}
