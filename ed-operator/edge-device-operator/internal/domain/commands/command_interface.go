package domain_commands

import (
	"context"
	"ed-operator/internal/domain/interfaces"
)

type Env map[string]string

type ICommand interface {
	SetContext(ctx context.Context)
	SetReconcilerClient(drc interfaces.IReconcilerClient)
	Execute() error
}
