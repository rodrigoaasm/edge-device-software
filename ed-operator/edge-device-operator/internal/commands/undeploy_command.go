package commands

import (
	"context"
	"ed-operator/internal/domain/interfaces"

	appsv1 "k8s.io/api/apps/v1"
	corekubev1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

	if dErr := d.reconcilerClient.Delete(d.ctx, dep); dErr != nil {
		return nil, dErr
	}
	log.Info("Searching for ports attached to the '" + d.Name + "' service")
	if dErr := d.deleteService(corekubev1.ServiceTypeClusterIP); dErr != nil {
		return nil, dErr
	}
	if dErr := d.deleteService(corekubev1.ServiceTypeNodePort); dErr != nil {
		return nil, dErr
	}

	return nil, nil
}

func (d *UndeployCommand) deleteService(t corekubev1.ServiceType) error {
	log := log.FromContext(d.ctx)

	var svc corev1.Service
	serviceClusterName := GetServiceName(d.Name, t)
	if dErr := d.reconcilerClient.Get(d.ctx, types.NamespacedName{
		Name:      serviceClusterName,
		Namespace: "ed-system",
	}, &svc); dErr == nil {
		log.Info("Removing the '" + serviceClusterName + "' service")
		if dErr := d.reconcilerClient.Delete(d.ctx, &svc); dErr != nil {
			return dErr
		}
	}

	return nil
}
