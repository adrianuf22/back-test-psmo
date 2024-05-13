package transaction

import "time"

type Service interface {
	Create(*Model) error
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
	eventDate       time.Time
}

func NewModel(accountId int64, operationTypeId int, amount float64) *Model {
	return &Model{
		accountID:       accountId,
		operationTypeID: OperationTypeID(operationTypeId),
		amount:          amount,
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
		return asPositiveBalance(m.amount)
	}

	return asNegativeBalance(m.amount)
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

func asPositiveBalance(amount float64) float64 {
	if amount < 0 {
		return amount * -1
	}

	return amount
}

func asNegativeBalance(amount float64) float64 {
	if amount > 0 {
		return amount * -1
	}

	return amount
}
