package undeploy

import (
	"context"
	"ed-operator/internal/domain/interfaces"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type UndeployCommand struct {
	CorrelationId string
	Name          string

	ctx              context.Context
	reconcilerClient interfaces.IReconcilerClient
}

func NewUndeployCommand(correlationId string, name string) *UndeployCommand {
	return &UndeployCommand{
		CorrelationId: correlationId,
		Name:          name,
	}
}

func (d *UndeployCommand) GetCorrelationId() string {
	return d.CorrelationId
}

func (d *UndeployCommand) SetContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *UndeployCommand) SetReconcilerClient(drc interfaces.IReconcilerClient) {
	d.reconcilerClient = drc
}

func (d *UndeployCommand) Execute() (interface{}, error) {
	log := log.FromContext(d.ctx)
	log.Info("Undeploying " + d.Name)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: "ed-system",
		},
	}
	dErr := d.reconcilerClient.Delete(d.ctx, dep)

	return nil, dErr
}
