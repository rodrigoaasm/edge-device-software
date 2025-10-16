package deploy

import (
	"context"
	cmd_commons "ed-operator/internal/domain/commands/commons"
	"ed-operator/internal/domain/entities"
	"ed-operator/internal/domain/interfaces"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	corekubev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DeployCommand struct {
	CorrelationId string
	Microservice  *entities.Microservice

	ctx              context.Context
	reconcilerClient interfaces.IReconcilerClient
}

func NewDeployCommand(
	correlationId string,
	microservice *entities.Microservice,
) *DeployCommand {
	return &DeployCommand{
		CorrelationId: correlationId,
		Microservice:  microservice,

		ctx:              context.Background(),
		reconcilerClient: nil,
	}
}

func (d *DeployCommand) GetCorrelationId() string {
	return d.CorrelationId
}

func (d *DeployCommand) SetContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *DeployCommand) SetReconcilerClient(drc interfaces.IReconcilerClient) {
	d.reconcilerClient = drc
}

func (d *DeployCommand) Execute() (interface{}, error) {
	log := log.FromContext(d.ctx)
	repl := int32(1)

	podSpec := d.createPodSpec()
	if d.Microservice.Port > 0 {
		podSpec.Containers[0].Ports = []corekubev1.ContainerPort{{ContainerPort: int32(d.Microservice.Port)}}
	}

	log.Info("Deploying " + d.Microservice.Name + "::" + d.Microservice.Image)
	if dErr := d.reconcilerClient.Create(d.ctx, d.createDeploymentSchema(podSpec, repl)); dErr != nil {
		return nil, dErr
	}

	log.Info("Open internal port (" + string(d.Microservice.InternalPort) + ") to " + d.Microservice.Name)
	if d.Microservice.InternalPort > 0 {
		if dErr := d.reconcilerClient.Create(d.ctx, d.createServiceSchema(apiv1.ServiceTypeClusterIP)); dErr != nil {
			return nil, dErr
		}
	}

	log.Info("Open external port (" + string(d.Microservice.ExternalPort) + ") to " + d.Microservice.Name)
	if d.Microservice.ExternalPort > 0 {
		if dErr := d.reconcilerClient.Create(d.ctx, d.createServiceSchema(apiv1.ServiceTypeNodePort)); dErr != nil {
			return nil, dErr
		}
	}

	return nil, nil
}

func (d *DeployCommand) createPodSpec() corekubev1.PodSpec {
	var envVars []corekubev1.EnvVar
	for key, value := range d.Microservice.Env {
		envVars = append(envVars, corekubev1.EnvVar{Name: key, Value: value})
	}

	podSpec := corekubev1.PodSpec{
		Containers: []corekubev1.Container{
			{
				Name:  d.Microservice.Name,
				Image: d.Microservice.Image,
				Env:   envVars,
				Resources: corekubev1.ResourceRequirements{
					Limits: corekubev1.ResourceList{
						corekubev1.ResourceCPU: *resource.NewMilliQuantity(
							int64(d.Microservice.LimitCPU), resource.DecimalSI,
						),
						corekubev1.ResourceMemory: *resource.NewQuantity(
							int64(d.Microservice.LimitMemory)*1024*1024, resource.BinarySI,
						),
					},
					Requests: corekubev1.ResourceList{
						corekubev1.ResourceCPU: *resource.NewMilliQuantity(
							int64(d.Microservice.RequestCPU), resource.DecimalSI,
						),
						corekubev1.ResourceMemory: *resource.NewQuantity(
							int64(d.Microservice.RequestMemory)*1024*1024, resource.BinarySI,
						),
					},
				},
			},
		},
	}

	if d.Microservice.Name == "operator" {
		podSpec.ServiceAccountName = "controller-manager"
	}

	if _, err := strconv.Atoi(d.Microservice.PriorityProfile); err != nil {
		podSpec.PriorityClassName = d.Microservice.PriorityProfile
	}

	return podSpec
}

func (d *DeployCommand) createDeploymentSchema(
	podSpec corekubev1.PodSpec,
	repl int32,
) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Microservice.Name,
			Namespace: "ed-system",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &repl,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": d.Microservice.Name},
			},
			Template: corekubev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": d.Microservice.Name},
				},
				Spec: podSpec,
			},
		},
	}
}

func (d *DeployCommand) createServiceSchema(t corekubev1.ServiceType) *apiv1.Service {
	portSchema := apiv1.ServicePort{
		Protocol:   apiv1.ProtocolTCP,
		Port:       int32(d.Microservice.Port),
		TargetPort: intstr.FromInt(int(d.Microservice.InternalPort)),
	}

	if t == corekubev1.ServiceTypeNodePort {
		portSchema.NodePort = int32(d.Microservice.ExternalPort)
	}

	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmd_commons.GetServiceName(d.Microservice.Name, t),
			Namespace: "ed-system",
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": d.Microservice.Name,
			},
			Type: t,
			Ports: []apiv1.ServicePort{
				portSchema,
			},
		},
	}
}
