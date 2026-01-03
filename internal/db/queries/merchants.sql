-- name: CreateMerchant :one
INSERT INTO merchants (name, email)
VALUES ($1, $2)
RETURNING id, name, email, status, created_at, updated_at;

-- name: GetMerchantByAPIKey :one
SELECT m.id, m.name, m.email, m.status, m.created_at, m.updated_at
FROM merchants m
JOIN merchant_api_keys mak ON m.id = mak.merchant_id
WHERE mak.api_key = $1 AND mak.status = 'active' AND m.status = 'active';

-- name: GetMerchant :one
SELECT id, name, email, status, created_at, updated_at
FROM merchants
WHERE id = $1;

-- name: CreateMerchantAPIKey :one
INSERT INTO merchant_api_keys (merchant_id, api_key, secret_key)
VALUES ($1, $2, $3)
RETURNING id, merchant_id, api_key, status, created_at, expires_at;
