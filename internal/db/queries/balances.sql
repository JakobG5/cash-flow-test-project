-- name: GetMerchantBalances :many
SELECT id, merchant_id, currency, available_balance, total_deposit, total_transaction_count, last_updated
FROM merchant_balances
WHERE merchant_id = $1;

-- name: GetMerchantBalance :one
SELECT id, merchant_id, currency, available_balance, total_deposit, total_transaction_count, last_updated
FROM merchant_balances
WHERE merchant_id = $1 AND currency = $2;

-- name: UpdateMerchantBalance :one
UPDATE merchant_balances
SET available_balance = $3, total_deposit = $4, total_transaction_count = $5, last_updated = NOW()
WHERE merchant_id = $1 AND currency = $2
RETURNING id, merchant_id, currency, available_balance, total_deposit, total_transaction_count, last_updated;

-- name: CreateMerchantBalance :one
INSERT INTO merchant_balances (merchant_id, currency)
VALUES ($1, $2)
RETURNING id, merchant_id, currency, available_balance, total_deposit, total_transaction_count, last_updated;

-- name: IncrementMerchantBalance :one
INSERT INTO merchant_balances (merchant_id, currency, available_balance, total_deposit, total_transaction_count)
VALUES ($1, $2, $3::decimal - $4::decimal, $3::decimal, 1)
ON CONFLICT (merchant_id, currency)
DO UPDATE SET
    available_balance = merchant_balances.available_balance + $3::decimal - $4::decimal,
    total_deposit = merchant_balances.total_deposit + $3::decimal,
    total_transaction_count = merchant_balances.total_transaction_count + 1,
    last_updated = NOW()
RETURNING id, merchant_id, currency, available_balance, total_deposit, total_transaction_count, last_updated;
