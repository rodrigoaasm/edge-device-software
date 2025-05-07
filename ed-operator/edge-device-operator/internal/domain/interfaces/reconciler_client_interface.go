package interfaces

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IReconcilerClient interface {
	Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error
	Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error
	Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error
}
