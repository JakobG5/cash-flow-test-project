-- name: CreatePaymentIntent :one
INSERT INTO payment_intents (payment_intent_id, merchant_id, amount, currency, description, callback_url, nonce, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, payment_intent_id, merchant_id, amount, currency, description, callback_url, nonce, status, metadata, created_at, updated_at, expires_at;

-- name: GetPaymentIntent :one
SELECT pi.id, pi.payment_intent_id, pi.merchant_id, pi.amount, pi.currency, pi.description, pi.callback_url, pi.nonce, pi.status, pi.metadata, pi.created_at, pi.updated_at, pi.expires_at
FROM payment_intents pi
WHERE pi.payment_intent_id = $1;

-- name: GetPaymentIntentByID :one
SELECT id, payment_intent_id, merchant_id, amount, currency, description, callback_url, nonce, status, metadata, created_at, updated_at, expires_at
FROM payment_intents
WHERE id = $1;

-- name: GetPaymentIntentByNonce :one
SELECT id, payment_intent_id, merchant_id, amount, currency, description, callback_url, nonce, status, metadata, created_at, updated_at, expires_at
FROM payment_intents
WHERE merchant_id = $1 AND nonce = $2;

-- name: UpdatePaymentIntentStatus :one
UPDATE payment_intents
SET status = $2, updated_at = NOW()
WHERE id = $1 AND status = $3
RETURNING id, payment_intent_id, merchant_id, amount, currency, description, callback_url, nonce, status, metadata, created_at, updated_at, expires_at;

-- name: LockPaymentIntentForProcessing :one
SELECT pi.id, pi.payment_intent_id, pi.merchant_id, pi.amount, pi.currency, pi.description, pi.callback_url, pi.nonce, pi.status, pi.metadata, pi.created_at, pi.updated_at, pi.expires_at
FROM payment_intents pi
WHERE pi.id = $1 AND pi.status = $2
FOR UPDATE;

-- name: CreatePaymentTransaction :one
INSERT INTO payment_transactions (payment_intent_id, merchant_id, amount, currency, payment_method, fee_amount, account_number)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, payment_intent_id, merchant_id, amount, currency, status, third_party_reference, payment_method, fee_amount, account_number, processed_at, created_at, updated_at;

-- name: UpdatePaymentTransactionStatus :one
UPDATE payment_transactions
SET status = $2, third_party_reference = $3, processed_at = NOW(), updated_at = NOW()
WHERE id = $1 AND status = $4
RETURNING id, payment_intent_id, merchant_id, amount, currency, status, third_party_reference, payment_method, fee_amount, account_number, processed_at, created_at, updated_at;

-- name: GetPaymentTransaction :one
SELECT id, payment_intent_id, merchant_id, amount, currency, status, third_party_reference, payment_method, fee_amount, account_number, processed_at, created_at, updated_at
FROM payment_transactions
WHERE id = $1;
