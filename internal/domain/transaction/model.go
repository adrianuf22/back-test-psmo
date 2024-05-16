package transaction

import (
	"errors"
	"time"
)

type Service interface {
	Create(*Model) error
	ReadAllPurchases(int64) ([]Model, error)
	Save([]Model) error
}

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
	return &Model{
		accountID:       accountId,
		operationTypeID: OperationTypeID(operationTypeId),
		amount:          amount,
		eventDate:       time.Now(),
	}
}

func NewModel(id int64, accountId int64, operationTypeId int, amount float64, balance float64, eventDate time.Time) *Model {
	return &Model{
		id:              id,
		accountID:       accountId,
		operationTypeID: OperationTypeID(operationTypeId),
		amount:          amount,
		balance:         balance,
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
	if m.operationTypeID == Payment {
		return asPositiveAmount(m.amount)
	}

	return asNegativeAmount(m.amount)
}

func (m *Model) Balance() float64 {
	return m.balance
}

func (m *Model) Discharge(amount float64) (float64, error) {
	if m.operationTypeID == Payment {
		return amount, errors.New("Operation type Payment is not allowed")
	}

	amount = amount - asPositiveAmount(m.balance)
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

func asPositiveAmount(amount float64) float64 {
	if amount < 0 {
		return amount * -1
	}

	return amount
}

func asNegativeAmount(amount float64) float64 {
	if amount > 0 {
		return amount * -1
	}

	return amount
}
