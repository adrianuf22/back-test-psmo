INSERT INTO transactions (
	account_id,
	operation_type,
	amount,
	balance,
	event_date
) VALUES ($1, $2, $3, $4, $5) RETURNING id;