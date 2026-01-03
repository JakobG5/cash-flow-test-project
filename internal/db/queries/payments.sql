-- name: CreatePaymentIntent :one
INSERT INTO payment_intents (transaction_id, merchant_id, amount, currency, reference, callback_url, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, transaction_id, merchant_id, amount, currency, reference, callback_url, status, metadata, created_at, updated_at, expires_at;

-- name: GetPaymentIntent :one
SELECT id, transaction_id, merchant_id, amount, currency, reference, callback_url, status, metadata, created_at, updated_at, expires_at
FROM payment_intents
WHERE transaction_id = $1;

-- name: GetPaymentIntentByID :one
SELECT id, transaction_id, merchant_id, amount, currency, reference, callback_url, status, metadata, created_at, updated_at, expires_at
FROM payment_intents
WHERE id = $1;

-- name: UpdatePaymentIntentStatus :one
UPDATE payment_intents
SET status = $2, updated_at = NOW()
WHERE id = $1 AND status = $3
RETURNING id, transaction_id, merchant_id, amount, currency, reference, callback_url, status, metadata, created_at, updated_at, expires_at;

-- name: CreatePaymentTransaction :one
INSERT INTO payment_transactions (transaction_id, payment_intent_id, merchant_id, amount, currency, payment_method)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, transaction_id, payment_intent_id, merchant_id, amount, currency, status, third_party_reference, payment_method, fee_amount, processed_at, created_at, updated_at;

-- name: UpdatePaymentTransactionStatus :one
UPDATE payment_transactions
SET status = $2, third_party_reference = $3, processed_at = NOW(), updated_at = NOW()
WHERE transaction_id = $1 AND status = $4
RETURNING id, transaction_id, payment_intent_id, merchant_id, amount, currency, status, third_party_reference, payment_method, fee_amount, processed_at, created_at, updated_at;

-- name: GetPaymentTransaction :one
SELECT id, transaction_id, payment_intent_id, merchant_id, amount, currency, status, third_party_reference, payment_method, fee_amount, processed_at, created_at, updated_at
FROM payment_transactions
WHERE transaction_id = $1;
