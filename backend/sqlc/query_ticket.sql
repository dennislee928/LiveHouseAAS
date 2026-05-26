-- name: GetTicket :one
SELECT id, order_id, ticket_type_id, event_id, code, qr_secret, status, used_at, created_at
FROM tickets
WHERE id = $1;

-- name: ListTicketsByOrder :many
SELECT id, order_id, ticket_type_id, event_id, code, qr_secret, status, used_at, created_at
FROM tickets
WHERE order_id = $1;

-- name: ListTicketsByUser :many
SELECT t.id, t.order_id, t.ticket_type_id, t.event_id, t.code, t.qr_secret, t.status, t.used_at, t.created_at
FROM tickets t
JOIN orders o ON t.order_id = o.id
WHERE o.user_id = $1
ORDER BY t.created_at DESC;

-- name: CreateTickets :copyfrom
INSERT INTO tickets (id, order_id, ticket_type_id, event_id, code, qr_secret, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6, 'active', NOW());

-- name: UseTicket :one
UPDATE tickets SET status = 'used', used_at = NOW()
WHERE id = $1 AND status = 'active'
RETURNING id, order_id, ticket_type_id, event_id, code, qr_secret, status, used_at, created_at;
