package health

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

func (u *Usecase) GetReadinessStatus() *Status {
	err := u.service.Readiness()
	if err != nil {
		return &OutOfService
	}

	return &UP
}
