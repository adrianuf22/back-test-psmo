package health

import "context"

type Status struct {
	Value string `json:"status"`
}

var (
	UP           = Status{"UP"}
	OutOfService = Status{"OUT_OF_SERVICE"}
)

type Service interface {
	Readiness(ctx context.Context) error
}
