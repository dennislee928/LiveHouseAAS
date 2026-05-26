-- name: ListTicketTypesByEvent :many
SELECT id, event_id, name, description, price, quantity, max_per_order, sale_start_at, sale_end_at, status, created_at, updated_at
FROM ticket_types
WHERE event_id = $1
ORDER BY price;

-- name: GetTicketType :one
SELECT id, event_id, name, description, price, quantity, max_per_order, sale_start_at, sale_end_at, status, created_at, updated_at
FROM ticket_types
WHERE id = $1;

-- name: CreateTicketType :one
INSERT INTO ticket_types (id, event_id, name, description, price, quantity, max_per_order, sale_start_at, sale_end_at, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'active', NOW(), NOW())
RETURNING id, event_id, name, description, price, quantity, max_per_order, sale_start_at, sale_end_at, status, created_at, updated_at;

-- name: UpdateTicketType :one
UPDATE ticket_types
SET name = $2, description = $3, price = $4, quantity = $5, max_per_order = $6, sale_start_at = $7, sale_end_at = $8, status = $9, updated_at = NOW()
WHERE id = $1
RETURNING id, event_id, name, description, price, quantity, max_per_order, sale_start_at, sale_end_at, status, created_at, updated_at;

-- name: DeleteTicketType :exec
DELETE FROM ticket_types WHERE id = $1;
