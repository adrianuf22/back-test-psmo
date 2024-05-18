package transaction

import (
	"errors"
	"time"
)

var (
	ErrOperationTypeNotAllowed = errors.New("operation type is not allowed")
	ErrInsufficientAmount      = errors.New("insufficient amount for operation")
)

type OperationTypeID int

const (
	_ OperationTypeID = iota
	CashPurchase
	InstallmentPurchase
	Withdrawal
	Payment
)

type Model struct {
	id              int64
	accountID       int64
	operationTypeID OperationTypeID
	amount          float64
	balance         float64
	eventDate       time.Time
}

func NewTransaction(accountId int64, operationTypeId int, amount float64) *Model {
	o := OperationTypeID(operationTypeId)

	return &Model{
		accountID:       accountId,
		operationTypeID: o,
		amount:          moneyByOperation(amount, o),
		balance:         moneyByOperation(amount, o),
		eventDate:       time.Now(),
	}
}

func NewModel(id int64, accountId int64, operationTypeId int, amount float64, balance float64, eventDate time.Time) *Model {
	o := OperationTypeID(operationTypeId)

	return &Model{
		id:              id,
		accountID:       accountId,
		operationTypeID: OperationTypeID(operationTypeId),
		amount:          moneyByOperation(amount, o),
		balance:         moneyByOperation(balance, o),
		eventDate:       time.Now(),
	}
}

func (m *Model) ID() int64 {
	return m.id
}

func (m *Model) SetID(id int64) {
	m.id = id
}

func (m *Model) AccountID() int64 {
	return m.accountID
}

func (m *Model) OperationTypeID() int {
	return int(m.operationTypeID)
}

func (m *Model) Amount() float64 {
	return m.amount
}

func (m *Model) Balance() float64 {
	return m.balance
}

func (m *Model) Discharge(amount float64) (float64, error) {
	if amount <= 0 {
		return amount, ErrInsufficientAmount
	}

	if m.operationTypeID == Payment {
		return amount, ErrOperationTypeNotAllowed
	}

	amount = amount - (-m.balance)
	if amount >= 0 {
		m.balance = 0
		return amount, nil
	}
	m.balance = amount

	return 0, nil
}

func (m *Model) EventDate() time.Time {
	return m.eventDate
}

func (m *Model) Validate() map[string]string {
	errs := make(map[string]string)

	if m.amount == 0 {
		errs["documentNumber"] = "document number is required"
	}

	return errs
}

func moneyByOperation(amount float64, operation OperationTypeID) float64 {
	if amount < 0 {
		amount = -amount
	}

	if operation == Payment {
		return amount
	}

	return -amount
}
