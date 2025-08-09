package controller

import (
	"context"
	"ed-operator/internal/app/mqtt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type PodPreemptionReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	MqttClient *mqtt.MQTTClient
}

func (r *PodPreemptionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	//logger := log.FromContext(context.Background())

	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Named("podpreemptionwatcher").
		WithEventFilter(predicate.Funcs{
			UpdateFunc: func(e event.UpdateEvent) bool {
				// oldPod := e.ObjectOld.(*corev1.Pod)
				// newPod := e.ObjectNew.(*corev1.Pod)

				return true // oldPod.Status.Phase != newPod.Status.Phase || oldPod.Status.Reason != newPod.Status.Reason
			},
		}).
		Complete(r)
}

func (r *PodPreemptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var pod corev1.Pod
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Pod status changed", "phase:", pod.Status.Phase, "reason:", pod.Status.Reason)
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Terminated != nil {
			state := cs.State.Terminated
			logger.Info("Container terminated", "container", cs.Name, "reason", state.Reason, "exitCode", state.ExitCode)
		}
	}

	return ctrl.Result{}, nil
}
