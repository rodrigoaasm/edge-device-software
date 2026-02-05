package commands

import (
	"context"
	"ed-operator/internal/domain/entities"
	"ed-operator/internal/domain/interfaces"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type ListDeployCommand struct {
	CorrelationId string

	reconcilerClient interfaces.IReconcilerClient
	ctx              context.Context
}

func NewListDeployCommand(correlationId string) *ListDeployCommand {
	return &ListDeployCommand{
		CorrelationId: correlationId,
	}
}

func (c *ListDeployCommand) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *ListDeployCommand) SetReconcilerClient(drc interfaces.IReconcilerClient) {
	c.reconcilerClient = drc
}

func (c *ListDeployCommand) GetCorrelationId() string {
	return c.CorrelationId
}

func (h *ListDeployCommand) Execute() (interface{}, error) {
	log := log.FromContext(h.ctx)
	log.Info("Searching deployments...")
	actualDeployments := &appsv1.DeploymentList{}
	if err := h.reconcilerClient.List(h.ctx, actualDeployments, client.InNamespace("ed-system")); err != nil {
		return nil, fmt.Errorf("unable to list deployments: %v", err)
	}

	var microserviceStatus []entities.MicroserviceSimpleStatus
	for _, actualDeployment := range actualDeployments.Items {
		status := actualDeployment.Status
		microserviceStatus = append(microserviceStatus, entities.MicroserviceSimpleStatus{
			Name:    actualDeployment.Name,
			Image:   actualDeployment.Spec.Template.Spec.Containers[0].Image,
			Healthy: status.Replicas == status.AvailableReplicas && status.AvailableReplicas > 0,
		})
	}

	log.Info("Deployments found.")
	return microserviceStatus, nil
}
