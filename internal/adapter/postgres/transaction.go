package postgres

import (
	"database/sql"
	_ "embed"
	"log/slog"

	"github.com/adrianuf22/back-test-psmo/internal/domain/transaction"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
)

//go:embed sql/create_transaction.sql
var createTransactionSql string

type transactRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) *transactRepo {
	return &transactRepo{
		db: db,
	}
}

func (r *transactRepo) Create(model *transaction.Model) error {
	var insertedId int64

	iAmount := int64(model.Amount() * 100)

	err := r.db.QueryRow(createTransactionSql, model.AccountID(), model.OperationTypeID(), iAmount, model.EventDate()).
		Scan(&insertedId)

	if err != nil {
		slog.Error(err.Error())
		return sentinel.WrapDBError(err)
	}

	model.SetID(insertedId)

	return nil
}
