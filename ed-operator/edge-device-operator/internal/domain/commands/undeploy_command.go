package domain_commands

import (
	"context"
	"ed-operator/internal/domain/interfaces"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type UndeployCommand struct {
	Name string

	ctx              context.Context
	reconcilerClient interfaces.IReconcilerClient
}

func NewUndeployCommand(name string) *UndeployCommand {
	return &UndeployCommand{
		Name: name,
	}
}

func (d *UndeployCommand) SetContext(ctx context.Context) {
	d.ctx = ctx
}

func (d *UndeployCommand) SetReconcilerClient(drc interfaces.IReconcilerClient) {
	d.reconcilerClient = drc
}

func (d *UndeployCommand) Execute() error {
	log := log.FromContext(d.ctx)
	log.Info("Undeploying " + d.Name)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: "ed-system",
		},
	}

	return d.reconcilerClient.Delete(d.ctx, dep)
}
