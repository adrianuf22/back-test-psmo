package transaction

import "time"

type Output struct {
	ID              int64           `json:"transaction_id"`
	AccountID       int64           `json:"account_id"`
	OperationTypeID OperationTypeID `json:"operation_type_id"`
	Amount          float64         `json:"amount"`
	EventDate       time.Time       `json:"event_date"`
}

func ToOutput(m *Model) *Output {
	return &Output{
		ID:              m.id,
		AccountID:       m.accountID,
		OperationTypeID: m.operationTypeID,
		Amount:          m.amount,
		EventDate:       m.eventDate,
	}
}
