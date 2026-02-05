package controllers

import (
	"context"
	command_commons "ed-operator/internal/commands"
	"ed-operator/internal/dtos"
	"ed-operator/internal/mqtt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type CriticalPodStatus string

const (
	CriticalPodStatusFailedScheduling CriticalPodStatus = "FailedScheduling"
	CriticalPodStatusStarted          CriticalPodStatus = "Started"
)

type CloudCommand string

const (
	CloudCommandPatch    CriticalPodStatus = "patch"
	CloudCommandDispatch CriticalPodStatus = "dispatch"
)

type PodPreemptionReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	cache      map[string]CriticalPodStatus
	MqttClient *mqtt.MQTTClient
}

func (r *PodPreemptionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	logger := log.FromContext(context.Background())
	start := time.Now()
	r.cache = make(map[string]CriticalPodStatus)

	return ctrl.NewControllerManagedBy(mgr).
		Named("podpreemptionwatcher").
		Watches(
			&corev1.Event{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				ev, ok := obj.(*corev1.Event)
				if ok && ev.CreationTimestamp.Time.After(start) && ev.InvolvedObject.Kind == "Pod" {
					if r.cache[ev.InvolvedObject.Name] == CriticalPodStatus(ev.Reason) {
						return nil
					}
					r.cache[ev.InvolvedObject.Name] = CriticalPodStatus(ev.Reason)

					if (ev.Reason == string(CriticalPodStatusFailedScheduling) && strings.Contains(ev.Message, "preemption")) ||
						ev.Reason == string(CriticalPodStatusStarted) {
						logger.Info("Pod:" + ev.InvolvedObject.Name + " reason:" + ev.Reason + " message:" + ev.Message)
						return []reconcile.Request{{
							NamespacedName: types.NamespacedName{
								Namespace: ev.InvolvedObject.Namespace,
								Name:      ev.InvolvedObject.Name,
							},
						}}
					}
				}
				return nil
			}),
		).Complete(r)
}

func (r *PodPreemptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(context.Background())

	if strings.Contains(req.NamespacedName.Name, "operator") {
		logger.Info("This event is about the Operator. Ignore it.")
		return ctrl.Result{}, nil
	}

	logger.Info("Getting more data about the pod:" + req.NamespacedName.Name + "that is " + string(r.cache[req.NamespacedName.Name]))
	var pod corev1.Pod
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Pod:" + pod.Name + " has been loaded. Publishing command to mqtt...")
	args := dtos.ArgsDto{
		Name:            pod.Name,
		Image:           pod.Spec.Containers[0].Image,
		PriorityProfile: pod.Spec.PriorityClassName,
		RequestMemory:   uint16(pod.Spec.Containers[0].Resources.Requests.Memory().Value()),
		RequestCPU:      uint16(pod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()),
		LimitMemory:     uint16(pod.Spec.Containers[0].Resources.Limits.Memory().Value()),
		LimitCPU:        uint16(pod.Spec.Containers[0].Resources.Limits.Cpu().MilliValue()),
		Env:             command_commons.EnvVarToEnvMap(pod.Spec.Containers[0].Env),
	}
	if len(pod.Spec.Containers[0].Ports) > 0 {
		args.Port = uint16(pod.Spec.Containers[0].Ports[0].ContainerPort)
	}

	cloudCmd := dtos.CommandDTO{
		Args:          args,
		CorrelationId: string(uuid.NewUUID()),
	}
	switch r.cache[req.NamespacedName.Name] {
	case CriticalPodStatusFailedScheduling:
		cloudCmd.Command = string(CloudCommandPatch)
	case CriticalPodStatusStarted:
		cloudCmd.Command = string(CloudCommandDispatch)
	}

	r.MqttClient.PublishCloudCommand(cloudCmd)
	return ctrl.Result{}, nil
}
