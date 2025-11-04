package promptports

import "context"

type UnitOfWorkPort interface {
	WithTransaction(ctx context.Context, work func(txCtx context.Context) (any, error)) (any, error)
}
