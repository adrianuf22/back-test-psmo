package postgres

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/adrianuf22/back-test-psmo/internal/domain/transaction"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/atomic"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
	"github.com/jackc/pgx/v5"
)

//go:embed sql/create_transaction.sql
var createTransactionSql string

//go:embed sql/read_transactions_by_op_type.sql
var readTransactionsByOpTypeSql string

//go:embed sql/update_transaction.sql
var updateTransactionSql string

type transactRepo struct {
	db pgx.Tx
}

type atomicTransactRepo struct {
	db *pgx.Conn
}

func NewTransactionRepo(db *pgx.Conn) atomicTransactRepo {
	return atomicTransactRepo{
		db: db,
	}
}

func (r atomicTransactRepo) Execute(
	ctx context.Context,
	op atomic.AtomicOperation[transaction.Repository],
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return sentinel.WrapDBError(err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	repoWithTransaction := transactRepo{db: tx}
	if err := op(ctx, repoWithTransaction); err != nil {
		return sentinel.WrapDBError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return sentinel.WrapDBError(err)
	}

	return nil
}

func (r transactRepo) ReadAllPurchases(ctx context.Context, accountId int64) ([]transaction.Model, error) {
	purchaseOps := []int32{int32(transaction.CashPurchase), int32(transaction.InstallmentPurchase)}
	rows, err := r.db.Query(ctx, readTransactionsByOpTypeSql, accountId, purchaseOps)
	if err != nil {
		slog.Error(err.Error())
		return nil, sentinel.WrapDBError(err)
	}
	defer rows.Close()

	var all []transaction.Model
	for rows.Next() {
		var (
			id              int64
			operationTypeID int
			amount,
			balance float64
			eventDate time.Time
		)

		err = rows.Scan(&id, &operationTypeID, &amount, &balance, &eventDate)
		if err != nil {
			slog.Error(err.Error())
			return nil, sentinel.WrapDBError(err)
		}

		all = append(all, *transaction.NewModel(id, accountId, operationTypeID, amount, balance, eventDate))
	}

	return all, nil
}

func (r transactRepo) Create(ctx context.Context, model *transaction.Model) error {
	var insertedId int64

	iAmount := int64(model.Amount() * 100)

	err := r.db.QueryRow(ctx, createTransactionSql, model.AccountID(), model.OperationTypeID(), iAmount, model.EventDate()).
		Scan(&insertedId)

	if err != nil {
		slog.Error(err.Error())
		return sentinel.WrapDBError(err)
	}

	model.SetID(insertedId)

	return nil
}

func (r transactRepo) UpdateAll(ctx context.Context, transactions []transaction.Model) error {
	if len(transactions) == 0 {
		return nil
	}

	b := &pgx.Batch{}
	for _, t := range transactions {
		b.Queue(updateTransactionSql, t.ID(), t.Balance())
	}

	results := r.db.SendBatch(ctx, b)
	_, err := results.Exec()
	if err != nil {
		return err
	}

	return nil
}
