-- name: ListSlotsByVenue :many
SELECT id, venue_id, date, start_time, end_time, status, created_at, updated_at
FROM slots
WHERE venue_id = $1
ORDER BY date, start_time;

-- name: ListAvailableSlotsByVenue :many
SELECT id, venue_id, date, start_time, end_time, status, created_at, updated_at
FROM slots
WHERE venue_id = $1 AND status = 'available' AND date >= CURRENT_DATE
ORDER BY date, start_time;

-- name: GetSlot :one
SELECT id, venue_id, date, start_time, end_time, status, created_at, updated_at
FROM slots
WHERE id = $1;

-- name: CreateSlot :one
INSERT INTO slots (id, venue_id, date, start_time, end_time, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
RETURNING id, venue_id, date, start_time, end_time, status, created_at, updated_at;

-- name: UpdateSlotStatus :one
UPDATE slots
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, venue_id, date, start_time, end_time, status, created_at, updated_at;

-- name: DeleteSlot :exec
DELETE FROM slots WHERE id = $1;

-- name: CheckSlotOverlap :one
SELECT COUNT(*) as cnt
FROM slots
WHERE venue_id = $1
  AND date = $2
  AND start_time < $4
  AND end_time > $3
  AND id != COALESCE($5, '00000000-0000-0000-0000-000000000000');
