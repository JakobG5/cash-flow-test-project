-- name: GetMerchantBalances :many
SELECT id, merchant_id, currency, available_balance, pending_balance, total_processed, last_updated
FROM merchant_balances
WHERE merchant_id = $1;

-- name: GetMerchantBalance :one
SELECT id, merchant_id, currency, available_balance, pending_balance, total_processed, last_updated
FROM merchant_balances
WHERE merchant_id = $1 AND currency = $2;

-- name: UpdateMerchantBalance :one
UPDATE merchant_balances
SET available_balance = $3, pending_balance = $4, total_processed = $5, last_updated = NOW()
WHERE merchant_id = $1 AND currency = $2
RETURNING id, merchant_id, currency, available_balance, pending_balance, total_processed, last_updated;

-- name: CreateMerchantBalance :one
INSERT INTO merchant_balances (merchant_id, currency)
VALUES ($1, $2)
RETURNING id, merchant_id, currency, available_balance, pending_balance, total_processed, last_updated;
