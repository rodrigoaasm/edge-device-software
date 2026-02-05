package commands

import (
	"context"
	"ed-operator/internal/domain/interfaces"
)

type ICommand interface {
	SetContext(ctx context.Context)
	SetReconcilerClient(drc interfaces.IReconcilerClient)
	GetCorrelationId() string
	Execute() (interface{}, error)
}
