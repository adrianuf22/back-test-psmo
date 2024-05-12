package domain

import "time"

type OperationTypeID int

const (
	_ OperationTypeID = iota
	CashPurchase
	InstallmentPurchase
	Withdrawal
	Payment
)

type Transaction struct {
	ID              int
	AccountID       int
	OperationTypeID OperationTypeID
	Amount          float32
	EventDate       time.Time
}

type TransactionService interface {
	Create(*Transaction) error
}
