package transaction

import (
	"context"

	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
)

type Usecase interface {
	CreateTransaction(context.Context, Input) (*Model, error)
}

type usecase struct {
	service        Service
	accountService account.Service
}

func NewUsecase(service Service, accountService account.Service) *usecase {
	return &usecase{
		service:        service,
		accountService: accountService,
	}
}

func (u *usecase) CreateTransaction(ctx context.Context, input Input) (*Model, error) {
	account, err := u.accountService.Read(ctx, input.AccountID)
	if err != nil {
		return nil, err
	}

	transaction := NewModel(account.ID(), input.OperationTypeID, input.Amount)
	err = u.service.Create(transaction)

	return transaction, err
}
