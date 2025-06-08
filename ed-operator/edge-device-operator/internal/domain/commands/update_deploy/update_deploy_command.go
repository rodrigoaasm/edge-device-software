package update_deploy

import (
	"context"
	"ed-operator/internal/domain/interfaces"

	appsv1 "k8s.io/api/apps/v1"
	corekubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Env map[string]string

type UpdateDeployCommand struct {
	CorrelationId string
	Name          string
	Image         string
	Env           Env

	ctx              context.Context
	reconcilerClient interfaces.IReconcilerClient
}

func NewUpdateDeployCommand(correlationId string, name string, image string, env Env) *UpdateDeployCommand {
	return &UpdateDeployCommand{
		CorrelationId: correlationId,
		Name:          name,
		Image:         image,
		Env:           env,
	}
}

func (d *UpdateDeployCommand) GetCorrelationId() string {
	return d.CorrelationId
}

func (d *UpdateDeployCommand) SetContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *UpdateDeployCommand) SetReconcilerClient(drc interfaces.IReconcilerClient) {
	d.reconcilerClient = drc
}

func (d *UpdateDeployCommand) Execute() (interface{}, error) {
	log := log.FromContext(d.ctx)
	log.Info("Update Deploying " + d.Name + "::" + d.Image)
	repl := int32(1)
	maxUnavailable := intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: 0,
	}
	maxSurge := intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: 1,
	}
	podSpec := corekubev1.PodSpec{
		Containers: []corekubev1.Container{
			{
				Name:  d.Name,
				Image: d.Image,
				Ports: []corekubev1.ContainerPort{{ContainerPort: 80}},
			},
		},
	}

	if d.Name == "operator" {
		podSpec.ServiceAccountName = "controller-manager"
	}

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
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       &maxSurge,
					MaxUnavailable: &maxUnavailable,
				},
			},
			Template: corekubev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": d.Name},
				},
				Spec: podSpec,
			},
		},
	}

	dErr := d.reconcilerClient.Update(d.ctx, dep)

	return nil, dErr
}
