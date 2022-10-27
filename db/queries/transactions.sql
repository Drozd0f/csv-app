-- name: CreateTransactions :copyfrom
INSERT INTO transactions (
    transaction_id, 
    request_id,
    terminal_id, partner_object_id,
    amount_total, amount_original,
    commission_ps, commission_client,
    commission_provider, date_input,
    date_post, status,
    payment_type, payment_number,
    service_id, service,
    payee_id, payee_name,
    payee_bank_mfo, payee_bank_account,
    payment_narrative
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
);
