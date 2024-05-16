BEGIN;
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    document_number varchar(11) NOT NULL,
    CONSTRAINT uniq_document_number UNIQUE(document_number)
);
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    account_id integer,
    operation_type integer,
    amount bigint,
    balance bigint,
    event_date timestamp,
    CONSTRAINT fk_account_id FOREIGN KEY(account_id) REFERENCES accounts(id)
);
COMMIT;

