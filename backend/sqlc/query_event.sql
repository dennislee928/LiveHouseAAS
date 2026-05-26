-- name: GetEvent :one
SELECT id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at
FROM events
WHERE id = $1;

-- name: ListEventsByVenue :many
SELECT id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at
FROM events
WHERE venue_id = $1
ORDER BY start_at DESC;

-- name: ListEventsByArtist :many
SELECT id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at
FROM events
WHERE artist_id = $1
ORDER BY start_at DESC;

-- name: CreateEvent :one
INSERT INTO events (id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'draft', NOW(), NOW())
RETURNING id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at;

-- name: UpdateEventStatus :one
UPDATE events
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at;
