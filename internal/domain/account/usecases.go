package account

import (
	"context"
)

type Usecase struct {
	service Service
}

func NewUsecase(service Service) *Usecase {
	return &Usecase{
		service: service,
	}
}

func (u *Usecase) GetAccountById(ctx context.Context, id int) (*Model, error) {
	return u.service.Read(ctx, id)
}

func (u *Usecase) CreateAccount(ctx context.Context, input Model) (*Model, error) {
	account := &Model{DocumentNumber: input.DocumentNumber}
	err := u.service.Create(ctx, account)

	return account, err
}
