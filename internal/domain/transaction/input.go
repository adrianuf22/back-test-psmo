package transaction

type Input struct {
	AccountID       int64   `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

func (i *Input) Validate() map[string]string {
	errs := make(map[string]string)

	if i.AccountID == 0 {
		errs["account_id"] = "account id is required"
	}

	if i.Amount == 0 {
		errs["amount"] = "amount is required"
	}

	if i.OperationTypeID <= 0 || i.OperationTypeID > int(Payment) {
		errs["operation_type_id"] = "invalid operation type id"
	}

	return errs
}
