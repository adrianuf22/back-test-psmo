package account

import "context"

type Repository interface {
	Read(context.Context, int64) (*Model, error)
	Create(context.Context, *Model) error
}
