package account

import (
	"context"
)

type Usecase interface {
	GetAccountById(context.Context, int64) (*Model, error)
	CreateAccount(context.Context, Input) (*Model, error)
}

type usecase struct {
	service Service
}

func NewUsecase(service Service) *usecase {
	return &usecase{
		service: service,
	}
}

func (u *usecase) GetAccountById(ctx context.Context, id int64) (*Model, error) {
	return u.service.Read(ctx, id)
}

func (u *usecase) CreateAccount(ctx context.Context, input Input) (*Model, error) {
	account := &Model{documentNumber: input.DocumentNumber}
	err := u.service.Create(ctx, account)

	return account, err
}
