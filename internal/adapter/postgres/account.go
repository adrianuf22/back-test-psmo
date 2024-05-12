package postgres

import (
	"context"
	"database/sql"

	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/error/db"
)

type accountRepo struct {
	DB *sql.DB
}

func NewAccount(db *sql.DB) *accountRepo {
	return &accountRepo{
		DB: db,
	}
}

func (a *accountRepo) Read(ctx context.Context, documentNumber int) (*account.Model, error) {
	account := &account.Model{}
	err := a.DB.QueryRowContext(ctx, `SELECT id, document_number FROM accounts WHERE document_number = $1`, documentNumber).Scan(
		&account.ID,
		&account.DocumentNumber,
	)
	if err != nil {
		return nil, db.WrapError(err)
	}

	return account, nil
}

func (a *accountRepo) Create(ctx context.Context, account *account.Model) error {
	var id int64
	err := a.DB.QueryRow(`INSERT INTO accounts (document_number) VALUES ($1) RETURNING id`, account.DocumentNumber).Scan(&id)
	if err != nil {
		return db.WrapError(err)
	}

	account.SetID(id)

	return nil
}
