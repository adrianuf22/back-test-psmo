package transaction

import (
	"context"
	"errors"

	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/atomic"
)

type Usecase interface {
	CreateTransaction(context.Context, Input) (*Model, error)
}

type usecase struct {
	service        atomic.AtomicRepository[Repository]
	accountService account.Repository
}

func NewUsecase(atomicService atomic.AtomicRepository[Repository], accountService account.Repository) *usecase {
	return &usecase{
		service:        atomicService,
		accountService: accountService,
	}
}

func (u *usecase) CreateTransaction(ctx context.Context, input Input) (*Model, error) {
	account, err := u.accountService.Read(ctx, input.AccountID)
	if err != nil {
		return nil, err
	}

	transaction := NewTransaction(account.ID(), input.OperationTypeID, input.Amount)

	err = u.service.Execute(ctx, func(ctx context.Context, r Repository) error {
		err = r.Create(ctx, transaction)
		if err != nil {
			return err
		}

		if transaction.operationTypeID != Payment {
			return nil
		}

		purchases, err := r.ReadAllPurchases(ctx, input.AccountID)
		if err != nil {
			return err
		}

		paymentBalance := transaction.amount
		for ix, p := range purchases {
			paymentBalance, err = p.Discharge(paymentBalance)
			if errors.Is(err, ErrInsufficientAmount) {
				break
			}

			purchases[ix] = p
		}

		return r.UpdateAll(ctx, purchases)
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
