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

-- name: SliceTransactions :many
SELECT *
FROM transactions
WHERE (@transaction_id::INTEGER = 0 OR transaction_id = @transaction_id)
AND (@status::TEXT = 'default' OR status = @status)
AND (@payment_type::TEXT = 'default' OR payment_type = @payment_type)
AND ((@date_post_from::TIMESTAMP = '0001-01-01' AND @date_post_to::TIMESTAMP = '0001-01-01') OR date_post BETWEEN @date_post_from AND @date_post_to)
AND (@transaction_id::TEXT = '' OR payment_narrative SIMILAR TO '%' || @payment_narrative || '%')
AND (cardinality(@terminal_id::INTEGER[]) = 0 OR terminal_id = ANY(@terminal_id));
