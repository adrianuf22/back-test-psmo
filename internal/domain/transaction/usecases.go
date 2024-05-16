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

	transaction := NewTransaction(account.ID(), input.OperationTypeID, input.Amount)
	err = u.service.Create(transaction)

	if transaction.operationTypeID == Payment {
		purchases, err := u.service.ReadAllPurchases(input.AccountID)
		if err != nil {
			// TODO
			return nil, err
		}

		paymentBalance := input.Amount
		for _, p := range purchases {
			if paymentBalance <= 0 {
				break
			}
			paymentBalance, _ = p.Discharge(paymentBalance)
		}

		err = u.service.Save(purchases)
	}

	return transaction, err
}
