package domain_commands

import (
	"context"
	"ed-operator/internal/domain/interfaces"

	appsv1 "k8s.io/api/apps/v1"
	corekubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type DeployCommand struct {
	Name  string
	Image string
	Env   Env

	ctx              context.Context
	reconcilerClient interfaces.IReconcilerClient
}

func NewDeployCommand(name string, image string, env Env) *DeployCommand {
	return &DeployCommand{
		Name:  name,
		Image: image,
		Env:   env,
	}
}

func (d *DeployCommand) SetContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *DeployCommand) SetReconcilerClient(drc interfaces.IReconcilerClient) {
	d.reconcilerClient = drc
}

func (d *DeployCommand) Execute() error {
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
						},
					},
				},
			},
		},
	}

	return d.reconcilerClient.Create(d.ctx, dep)
}
