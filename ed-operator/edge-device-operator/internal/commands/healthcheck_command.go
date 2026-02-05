package commands

import (
	"context"
	"ed-operator/internal/domain/interfaces"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type HealthCheckCommand struct {
	CorrelationId string

	reconcilerClient interfaces.IReconcilerClient
	ctx              context.Context
}

func NewHealthCheckCommand(correlationId string) *HealthCheckCommand {
	return &HealthCheckCommand{
		CorrelationId: correlationId,
	}
}

func (c *HealthCheckCommand) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *HealthCheckCommand) SetReconcilerClient(drc interfaces.IReconcilerClient) {
	c.reconcilerClient = drc
}

func (c *HealthCheckCommand) GetCorrelationId() string {
	return c.CorrelationId
}

func (h *HealthCheckCommand) Execute() (interface{}, error) {
	log := log.FromContext(h.ctx)
	log.Info("Searching ed-operator deployment...")
	var actualDeployment appsv1.Deployment
	if err := h.reconcilerClient.Get(h.ctx, types.NamespacedName{
		Name:      "operator",
		Namespace: "ed-system",
	}, &actualDeployment); err != nil {
		return nil, fmt.Errorf("unable to get ed-operator deployment: %v", err)
	}

	return actualDeployment.Spec.Template.Spec.Containers[0].Image, nil
}
