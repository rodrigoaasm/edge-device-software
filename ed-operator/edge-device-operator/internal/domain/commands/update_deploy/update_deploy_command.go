package update_deploy

import (
	"context"
	command_commons "ed-operator/internal/domain/commands/commons"
	"ed-operator/internal/domain/entities"
	"ed-operator/internal/domain/interfaces"
	"errors"

	appsv1 "k8s.io/api/apps/v1"
	corekubev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Env map[string]string

type UpdateDeployCommand struct {
	CorrelationId string
	Microservice  *entities.Microservice

	ctx              context.Context
	reconcilerClient interfaces.IReconcilerClient
}

func NewUpdateDeployCommand(correlationId string, microservice *entities.Microservice) *UpdateDeployCommand {
	return &UpdateDeployCommand{
		CorrelationId: correlationId,
		Microservice:  microservice,
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
	log.Info("Update Deploying " + d.Microservice.Name + "::" + d.Microservice.Image)
	maxUnavailable := intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: 0,
	}
	maxSurge := intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: 1,
	}

	log.Info("Searching " + d.Microservice.Name + " deployment...")
	var actualDeployment appsv1.Deployment
	d.reconcilerClient.Get(d.ctx, types.NamespacedName{
		Name:      d.Microservice.Name,
		Namespace: "ed-system",
	}, &actualDeployment)
	log.Info(d.Microservice.Name + " deployment found.")

	log.Info("Updating " + d.Microservice.Name + " deployment...")
	if d.Microservice.Image != "" {
		actualDeployment.Spec.Template.Spec.Containers[0].Image = d.Microservice.Image
	}

	if d.Microservice.Env != nil {
		envVars := command_commons.EnvMapToEnvVar(d.Microservice.Env)
		var deviceIDEnv string
		if d.Microservice.Name == "operator" {
			for _, item := range actualDeployment.Spec.Template.Spec.Containers[0].Env {
				if item.Name == "DEVICE_ID" {
					deviceIDEnv = item.Value
					break
				}
			}
			if deviceIDEnv == "" {
				return nil, errors.New("device id not found in operator deployment")
			}

			envVars = append(envVars, corekubev1.EnvVar{Name: "DEVICE_ID", Value: deviceIDEnv})
		}
		actualDeployment.Spec.Template.Spec.Containers[0].Env = envVars
	}

	if d.Microservice.Name == "operator" {
		actualDeployment.Spec.Template.Spec.ServiceAccountName = "controller-manager"
	}

	if d.Microservice.Port > 0 {
		actualDeployment.Spec.Template.Spec.Containers[0].Ports = []corekubev1.ContainerPort{{ContainerPort: int32(d.Microservice.Port)}}
	}

	actualDeployment.Spec.Strategy = appsv1.DeploymentStrategy{
		Type: appsv1.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxSurge:       &maxSurge,
			MaxUnavailable: &maxUnavailable,
		},
	}

	log.Info("Applying changes in " + d.Microservice.Name + " deployment...")
	dErr := d.reconcilerClient.Update(d.ctx, &actualDeployment)
	log.Info("Changes applied in " + d.Microservice.Name + " deployment. Deployment updated")

	return nil, dErr
}
