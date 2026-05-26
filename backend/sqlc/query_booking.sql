-- name: GetBooking :one
SELECT id, slot_id, venue_id, artist_id, message, status, created_at, updated_at
FROM booking_requests
WHERE id = $1;

-- name: ListBookingsByVenue :many
SELECT id, slot_id, venue_id, artist_id, message, status, created_at, updated_at
FROM booking_requests
WHERE venue_id = $1
ORDER BY created_at DESC;

-- name: ListBookingsByArtist :many
SELECT id, slot_id, venue_id, artist_id, message, status, created_at, updated_at
FROM booking_requests
WHERE artist_id = $1
ORDER BY created_at DESC;

-- name: ListBookingsByStatus :many
SELECT id, slot_id, venue_id, artist_id, message, status, created_at, updated_at
FROM booking_requests
WHERE status = $1
ORDER BY created_at DESC;

-- name: CreateBooking :one
INSERT INTO booking_requests (id, slot_id, venue_id, artist_id, message, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, 'pending', NOW(), NOW())
RETURNING id, slot_id, venue_id, artist_id, message, status, created_at, updated_at;

-- name: UpdateBookingStatus :one
UPDATE booking_requests
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, slot_id, venue_id, artist_id, message, status, created_at, updated_at;
