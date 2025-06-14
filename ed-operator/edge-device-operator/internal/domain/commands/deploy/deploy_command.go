package deploy

import (
	"context"
	"ed-operator/internal/domain/interfaces"

	appsv1 "k8s.io/api/apps/v1"
	corekubev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Env map[string]string

type DeployCommand struct {
	CorrelationId string
	Name          string
	Image         string
	Env           Env

	ctx              context.Context
	reconcilerClient interfaces.IReconcilerClient
}

func NewDeployCommand(correlationId string, name string, image string, env Env) *DeployCommand {
	return &DeployCommand{
		CorrelationId: correlationId,
		Name:          name,
		Image:         image,
		Env:           env,
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
	log.Info("Deploying " + d.Name + "::" + d.Image)
	repl := int32(1)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: "ed-system",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &repl,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": d.Name},
			},
			Template: corekubev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": d.Name},
				},
				Spec: corekubev1.PodSpec{
					Containers: []corekubev1.Container{
						{
							Name:  d.Name,
							Image: d.Image,
							Ports: []corekubev1.ContainerPort{{ContainerPort: 80}},
							Resources: corekubev1.ResourceRequirements{
								Limits: corekubev1.ResourceList{
									corekubev1.ResourceCPU:    *resource.NewMilliQuantity(500, resource.DecimalSI),
									corekubev1.ResourceMemory: *resource.NewQuantity(256*1024*1024, resource.BinarySI),
								},
								Requests: corekubev1.ResourceList{
									corekubev1.ResourceCPU:    *resource.NewMilliQuantity(200, resource.DecimalSI),
									corekubev1.ResourceMemory: *resource.NewQuantity(128*1024*1024, resource.BinarySI),
								},
							},
						},
					},
				},
			},
		},
	}

	dErr := d.reconcilerClient.Create(d.ctx, dep)

	return nil, dErr
}
