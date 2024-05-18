package health

import "context"

type Usecase struct {
	service Service
}

func NewUsecase(service Service) *Usecase {
	return &Usecase{
		service: service,
	}
}

func (u *Usecase) GetLivenessStatus() *Status {
	return &UP
}

func (u *Usecase) GetReadinessStatus(ctx context.Context) *Status {
	err := u.service.Readiness(ctx)
	if err != nil {
		return &OutOfService
	}

	return &UP
}
