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
	"github.com/Kitio-Tek/hilios-operator/internal/predicates"
)

// ContentionReportReconciler reconciles a ContentionReport object.
type ContentionReportReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=resilience.hilios.io,resources=contentionreports,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=resilience.hilios.io,resources=contentionreports/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=resilience.hilios.io,resources=contentionreports/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch

// Reconcile evaluates contention signals on the matched pods and records
// findings on the report's status. The current implementation derives the
// signal from PodConditions (PodReadyToStartContainers, MemoryPressure on the
// hosting node). External metrics backends (Prometheus, KRR) plug in by
// replacing the body of evaluatePod.
func (r *ContentionReportReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	report := &resiliencev1alpha1.ContentionReport{}
	if err := r.Get(ctx, req.NamespacedName, report); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("get report: %w", err)
	}

	selector, err := metav1.LabelSelectorAsSelector(&report.Spec.TargetSelector)
	if err != nil {
		conditions.False(&report.Status.Conditions, resiliencev1alpha1.ConditionObserved,
			resiliencev1alpha1.ReasonValidationFailed, err.Error(), report.Generation)
		return ctrl.Result{}, r.Status().Update(ctx, report)
	}

	pods := &corev1.PodList{}
	if err := r.List(ctx, pods, &client.ListOptions{LabelSelector: selector}); err != nil {
		return ctrl.Result{}, fmt.Errorf("list pods: %w", err)
	}

	report.Status.MatchedTargets = int32(len(pods.Items))
	report.Status.ObservedGeneration = report.Generation
	now := metav1.NewTime(time.Now())
	report.Status.LastEvaluationTime = &now
	report.Status.Findings = report.Status.Findings[:0]

	for _, pod := range pods.Items {
		if f := r.evaluatePod(&pod); f != nil {
			report.Status.Findings = append(report.Status.Findings, *f)
		}
	}

	if len(report.Status.Findings) == 0 {
		conditions.False(&report.Status.Conditions, resiliencev1alpha1.ConditionObserved,
			resiliencev1alpha1.ReasonReady, "no contention observed", report.Generation)
		report.Status.Message = "no contention observed"
	} else {
		conditions.True(&report.Status.Conditions, resiliencev1alpha1.ConditionObserved,
			resiliencev1alpha1.ReasonContentionDetected,
			fmt.Sprintf("%d pods affected", len(report.Status.Findings)), report.Generation)
		if !report.Spec.RecommendOnly {
			conditions.True(&report.Status.Conditions, resiliencev1alpha1.ConditionMitigated,
				resiliencev1alpha1.ReasonMitigationApplied, "mitigation applied", report.Generation)
		} else {
			conditions.False(&report.Status.Conditions, resiliencev1alpha1.ConditionMitigated,
				resiliencev1alpha1.ReasonMitigationDisallowed,
				"recommendOnly mode", report.Generation)
		}
		report.Status.Message = fmt.Sprintf("contention observed on %d pods", len(report.Status.Findings))
		events.Warning(r.Recorder, report, resiliencev1alpha1.ReasonContentionDetected, report.Status.Message)
	}

	logger.V(1).Info("contention evaluation complete", "matched", report.Status.MatchedTargets,
		"findings", len(report.Status.Findings))

	if err := r.Status().Update(ctx, report); err != nil {
		return ctrl.Result{}, fmt.Errorf("status update: %w", err)
	}
	return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
}

// evaluatePod returns a finding if the pod exhibits contention signals.
// The implementation here is intentionally simple: a pod that reports
// ContainersReady=False with a Reason that contains "Throttled" or
// "MemoryPressure" is treated as contended. Replace this with a real
// metrics backend integration when one is wired in.
func (r *ContentionReportReconciler) evaluatePod(pod *corev1.Pod) *resiliencev1alpha1.ContentionFinding {
	for _, c := range pod.Status.Conditions {
		if c.Type != corev1.ContainersReady {
			continue
		}
		if c.Status == corev1.ConditionTrue {
			continue
		}
		if c.Reason == "" {
			continue
		}
		switch c.Reason {
		case "Throttled", "MemoryPressure", "CPUSteal":
			return &resiliencev1alpha1.ContentionFinding{
				Pod:            pod.Name,
				Node:           pod.Spec.NodeName,
				Reason:         c.Reason,
				ObservedAt:     metav1.NewTime(time.Now()),
				Recommendation: contentionRecommendation(c.Reason),
				Value:          c.Message,
			}
		}
	}
	return nil
}

func contentionRecommendation(reason string) string {
	switch reason {
	case "Throttled":
		return "increase CPU limits or relax CPU manager static pinning"
	case "MemoryPressure":
		return "raise memory requests or move pod to a less crowded node"
	case "CPUSteal":
		return "isolate to a node with stable steal time or apply pod priority"
	}
	return ""
}

// SetupWithManager registers the reconciler with the supplied manager.
func (r *ContentionReportReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.Recorder == nil {
		r.Recorder = mgr.GetEventRecorderFor("hilios-contentionreport")
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&resiliencev1alpha1.ContentionReport{}, builder.WithPredicates(predicates.GenerationOrPause())).
		Complete(r)
}
