package account

type Output struct {
	ID             int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

func ToOutput(m *Model) *Output {
	return &Output{
		ID:             m.id,
		DocumentNumber: m.documentNumber,
	}
}
