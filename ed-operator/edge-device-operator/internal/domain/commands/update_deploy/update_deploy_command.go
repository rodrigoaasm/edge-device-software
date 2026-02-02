package update_deploy

import (
	"context"
	command_commons "ed-operator/internal/domain/commands/commons"
	"ed-operator/internal/domain/entities"
	"ed-operator/internal/domain/interfaces"
	"errors"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corekubev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
	err := d.reconcilerClient.Get(d.ctx, types.NamespacedName{
		Name:      d.Microservice.Name,
		Namespace: "ed-system",
	}, &actualDeployment)
	if err != nil {
		return nil, errors.New("unable to get " + d.Microservice.Name + " deployment: " + err.Error())
	}
	log.Info(d.Microservice.Name + " deployment found.")

	log.Info("Updating " + d.Microservice.Name + " deployment...")
	if d.Microservice.Image != "" {
		actualDeployment.Spec.Template.Spec.Containers[0].Image = d.Microservice.Image
	}

	depTemplateSpec := actualDeployment.Spec.Template.Spec
	if d.Microservice.Env != nil {
		log.Info("Update Envs for " + d.Microservice.Name)
		envVars := command_commons.EnvMapToEnvVar(d.Microservice.Env)
		var deviceIDEnv string
		if d.Microservice.Name == "operator" {
			log.Info("This update is for operator, update Env DEVICE_ID for " + d.Microservice.Name)
			for _, item := range depTemplateSpec.Containers[0].Env {
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
		depTemplateSpec.Containers[0].Env = envVars
	}

	if d.Microservice.Name == "operator" {
		depTemplateSpec.ServiceAccountName = "controller-manager"
	}

	if d.Microservice.Port > 0 {
		log.Info("Update port for" + d.Microservice.Name + ". From:" + depTemplateSpec.Containers[0].Ports[0].String() + " to:" + strconv.FormatUint(uint64(d.Microservice.Port), 10))
		depTemplateSpec.Containers[0].Ports = []corekubev1.ContainerPort{{ContainerPort: int32(d.Microservice.Port)}}
	}

	r := depTemplateSpec.Containers[0].Resources
	var requestCPU resource.Quantity
	var requestMemo resource.Quantity
	var limitCPU resource.Quantity
	var limitMemo resource.Quantity
	if d.Microservice.RequestCPU > 0 {
		log.Info("Update requestCPU fo " + d.Microservice.Name + ". From:" + r.Requests.Cpu().String() + " to:" + strconv.FormatUint(uint64(d.Microservice.RequestCPU), 10))
		requestCPU = *resource.NewMilliQuantity(
			int64(d.Microservice.RequestCPU), resource.DecimalSI,
		)
	} else {
		requestCPU = *r.Requests.Cpu()
	}

	if d.Microservice.RequestMemory > 0 {
		log.Info("Update requestMemory for " + d.Microservice.Name + ". From:" + r.Requests.Memory().String() + " to:" + strconv.FormatUint(uint64(d.Microservice.RequestMemory), 10))
		requestMemo = *resource.NewQuantity(
			int64(d.Microservice.RequestMemory)*1024*1024, resource.BinarySI,
		)
	} else {
		requestMemo = *r.Requests.Memory()
	}

	if d.Microservice.LimitCPU > 0 {
		log.Info("Update limitCPU for " + d.Microservice.Name + ". From:" + r.Limits.Cpu().String() + " to:" + strconv.FormatUint(uint64(d.Microservice.LimitCPU), 10))
		limitCPU = *resource.NewMilliQuantity(
			int64(d.Microservice.LimitCPU), resource.DecimalSI,
		)
	} else {
		limitCPU = *r.Limits.Cpu()
	}

	if d.Microservice.LimitMemory > 0 {
		log.Info("Update LimitMemory for " + d.Microservice.Name + ". From:" + r.Limits.Memory().String() + " to:" + strconv.FormatUint(uint64(d.Microservice.LimitMemory), 10))
		limitMemo = *resource.NewQuantity(
			int64(d.Microservice.LimitMemory)*1024*1024, resource.BinarySI,
		)
	} else {
		limitMemo = *r.Limits.Memory()
	}

	if d.Microservice.RequestCPU > 0 || d.Microservice.LimitCPU > 0 || d.Microservice.RequestMemory > 0 || d.Microservice.LimitMemory > 0 {
		depTemplateSpec.Containers[0].Resources = corekubev1.ResourceRequirements{
			Limits: corekubev1.ResourceList{
				corekubev1.ResourceCPU:    limitCPU,
				corekubev1.ResourceMemory: limitMemo,
			},
			Requests: corekubev1.ResourceList{
				corekubev1.ResourceCPU:    requestCPU,
				corekubev1.ResourceMemory: requestMemo,
			},
		}
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
