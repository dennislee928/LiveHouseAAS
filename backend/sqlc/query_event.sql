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

-- name: ListPublishedEvents :many
SELECT e.id, e.title, e.description, e.venue_id, e.artist_id, e.booking_id, e.start_at, e.end_at, e.status, e.created_at, e.updated_at,
       v.name AS venue_name, v.city AS venue_city,
       u.name AS artist_name
FROM events e
JOIN venues v ON e.venue_id = v.id
JOIN users u ON e.artist_id = u.id
WHERE e.status = 'published' AND e.start_at > NOW()
ORDER BY e.start_at;

-- name: CreateEvent :one
INSERT INTO events (id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'draft', NOW(), NOW())
RETURNING id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at;

-- name: UpdateEvent :one
UPDATE events
SET title = $2, description = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at;

-- name: UpdateEventStatus :one
UPDATE events
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, title, description, venue_id, artist_id, booking_id, start_at, end_at, status, created_at, updated_at;
