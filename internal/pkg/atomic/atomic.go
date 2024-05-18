package atomic

import "context"

type AtomicOperation[T any] func(context.Context, T) error

type AtomicRepository[T any] interface {
	Execute(context.Context, AtomicOperation[T]) error
}
