package postgres

import (
	"context"
	_ "embed"

	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
	"github.com/jackc/pgx/v5"
)

//go:embed sql/get_account.sql
var getAccountSql string

//go:embed sql/create_account.sql
var createAccountSql string

type accountRepo struct {
	db *pgx.Conn
}

func NewAccountRepo(db *pgx.Conn) *accountRepo {
	return &accountRepo{
		db: db,
	}
}

func (r *accountRepo) Read(ctx context.Context, id int64) (*account.Model, error) {
	var accountId int
	var documentNumber string

	err := r.db.QueryRow(ctx, getAccountSql, id).Scan(
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
	err := r.db.QueryRow(ctx, createAccountSql, model.DocumentNumber()).
		Scan(&insertedId)

	if err != nil {
		return sentinel.WrapDBError(err)
	}

	model.SetID(insertedId)

	return nil
}
