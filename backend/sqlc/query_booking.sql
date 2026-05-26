-- name: ListBookingsByVenue :many
SELECT br.id, br.slot_id, br.venue_id, br.artist_id, br.message, br.status, br.created_at, br.updated_at,
       s.date, s.start_time, s.end_time,
       u.name AS artist_name, u.email AS artist_email
FROM booking_requests br
JOIN slots s ON br.slot_id = s.id
JOIN users u ON br.artist_id = u.id
WHERE br.venue_id = $1
ORDER BY br.created_at DESC;

-- name: ListBookingsByArtist :many
SELECT br.id, br.slot_id, br.venue_id, br.artist_id, br.message, br.status, br.created_at, br.updated_at,
       s.date, s.start_time, s.end_time,
       v.name AS venue_name, v.city AS venue_city
FROM booking_requests br
JOIN slots s ON br.slot_id = s.id
JOIN venues v ON br.venue_id = v.id
WHERE br.artist_id = $1
ORDER BY br.created_at DESC;

-- name: GetBooking :one
SELECT id, slot_id, venue_id, artist_id, message, status, created_at, updated_at
FROM booking_requests
WHERE id = $1;

-- name: CreateBooking :one
INSERT INTO booking_requests (id, slot_id, venue_id, artist_id, message, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, 'pending', NOW(), NOW())
RETURNING id, slot_id, venue_id, artist_id, message, status, created_at, updated_at;

-- name: UpdateBookingStatus :one
UPDATE booking_requests
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, slot_id, venue_id, artist_id, message, status, created_at, updated_at;
