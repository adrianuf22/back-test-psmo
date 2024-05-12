package account

import "context"

type Model struct {
	ID             int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type Service interface {
	Read(context.Context, int) (*Model, error)
	Create(context.Context, *Model) error
}

func (m *Model) SetID(id int64) {
	m.ID = id
}

func (m *Model) Validate() map[string]string {
	errs := make(map[string]string)

	if m.DocumentNumber == "" {
		errs["documentNumber"] = "document number is required"
	}

	return errs
}
