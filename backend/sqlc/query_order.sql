-- name: GetOrder :one
SELECT id, user_id, event_id, total_amount, status, payment_method, paid_at, created_at, updated_at
FROM orders
WHERE id = $1;

-- name: ListOrdersByUser :many
SELECT id, user_id, event_id, total_amount, status, payment_method, paid_at, created_at, updated_at
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateOrder :one
INSERT INTO orders (id, user_id, event_id, total_amount, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, 'pending', NOW(), NOW())
RETURNING id, user_id, event_id, total_amount, status, payment_method, paid_at, created_at, updated_at;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, payment_method = COALESCE($3, payment_method), paid_at = CASE WHEN $2 = 'paid' THEN NOW() ELSE paid_at END, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, event_id, total_amount, status, payment_method, paid_at, created_at, updated_at;
