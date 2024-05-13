package account

type Input struct {
	DocumentNumber string `json:"document_number"`
}

func (i *Input) Validate() map[string]string {
	errs := make(map[string]string)

	if i.DocumentNumber == "" {
		errs["documentNumber"] = "document number is required"
	}

	return errs
}
