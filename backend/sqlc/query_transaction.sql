-- name: CreateTransaction :one
INSERT INTO transactions (id, order_id, provider, amount, fee, provider_tx_id, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, 'pending', NOW(), NOW())
RETURNING id, order_id, provider, amount, fee, provider_tx_id, status, created_at, updated_at;

-- name: UpdateTransactionStatus :one
UPDATE transactions
SET status = $2, provider_tx_id = COALESCE($3, provider_tx_id), updated_at = NOW()
WHERE id = $1
RETURNING id, order_id, provider, amount, fee, provider_tx_id, status, created_at, updated_at;
