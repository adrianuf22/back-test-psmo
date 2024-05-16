package postgres

import (
	"database/sql"
	_ "embed"
	"log/slog"
	"time"

	"github.com/adrianuf22/back-test-psmo/internal/domain/transaction"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
)

//go:embed sql/create_transaction.sql
var createTransactionSql string

//go:embed sql/read_all_purchases.sql
var readAllPurchasesSql string

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

func (r *transactRepo) ReadAllPurchases(accountId int64) (*[]transaction.Model, error) {
	rows, err := r.db.Query(readAllPurchasesSql, accountId)
	if err != nil {
		slog.Error(err.Error())
		return nil, sentinel.WrapDBError(err)
	}

	var all []transaction.Model
	for rows.Next() {
		var (
			id              int64
			operationTypeID int
			amount,
			balance float64
			eventDate time.Time
		)

		rows.Scan(&id, &operationTypeID, &amount, &balance, &eventDate)

		all = append(all, *transaction.NewModel(id, accountId, operationTypeID, amount, balance, eventDate))
	}
}

// Save([]Model) error
