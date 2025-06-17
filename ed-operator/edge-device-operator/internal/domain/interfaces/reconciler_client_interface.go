package interfaces

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IReconcilerClient interface {
	Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error
	List(ctx context.Context, obj client.ObjectList, opts ...client.ListOption) error
	Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error
	Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error
	Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error
}
