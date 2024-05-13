INSERT INTO transactions (
	account_id,
	operation_type,
	amount,
	event_date
) VALUES ($1, $2, $3, $4) RETURNING id;