package transaction

import "context"

type Repository interface {
	Create(context.Context, *Model) error
	ReadAllPurchases(context.Context, int64) ([]Model, error)
	UpdateAll(context.Context, []Model) error
}
