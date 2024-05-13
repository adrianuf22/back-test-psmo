package postgres

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
)

//go:embed sql/get_account.sql
var getAccountSql string

//go:embed sql/create_account.sql
var createAccountSql string

type accountRepo struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) *accountRepo {
	return &accountRepo{
		db: db,
	}
}

func (r *accountRepo) Read(ctx context.Context, id int64) (*account.Model, error) {
	var accountId int
	var documentNumber string

	err := r.db.QueryRowContext(ctx, getAccountSql, id).Scan(
		&accountId,
		&documentNumber,
	)
	if err != nil {
		return nil, sentinel.WrapDBError(err)
	}

	return account.NewModel(int64(accountId), documentNumber), nil
}

func (r *accountRepo) Create(ctx context.Context, model *account.Model) error {
	var insertedId int64
	err := r.db.QueryRow(createAccountSql, model.DocumentNumber()).
		Scan(&insertedId)

	if err != nil {
		return sentinel.WrapDBError(err)
	}

	model.SetID(insertedId)

	return nil
}
