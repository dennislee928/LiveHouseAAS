-- name: ListEventsByVenue :many
SELECT e.id, e.title, e.description, e.venue_id, e.artist_id, e.booking_id, e.start_at, e.end_at, e.status, e.created_at, e.updated_at,
       u.name AS artist_name
FROM events e
JOIN users u ON e.artist_id = u.id
WHERE e.venue_id = $1
ORDER BY e.start_at DESC;

-- name: ListEventsByArtist :many
SELECT e.id, e.title, e.description, e.venue_id, e.artist_id, e.booking_id, e.start_at, e.end_at, e.status, e.created_at, e.updated_at,
       v.name AS venue_name
FROM events e
JOIN venues v ON e.venue_id = v.id
WHERE e.artist_id = $1
ORDER BY e.start_at DESC;

-- name: GetEvent :one
SELECT id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at
FROM events
WHERE id = $1;

-- name: CreateEvent :one
INSERT INTO events (id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'draft', NOW(), NOW())
RETURNING id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at;

-- name: UpdateEventStatus :one
UPDATE events
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at;
