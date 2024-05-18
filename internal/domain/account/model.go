package account

type Model struct {
	id             int64
	documentNumber string
}

func NewModel(id int64, documentNumber string) *Model {
	return &Model{
		id:             id,
		documentNumber: documentNumber,
	}
}

func (m *Model) ID() int64 {
	return m.id
}

func (m *Model) SetID(id int64) {
	m.id = id
}

func (m *Model) DocumentNumber() string {
	return m.documentNumber
}
