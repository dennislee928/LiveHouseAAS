-- name: ListVenues :many
SELECT id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at
FROM venues
ORDER BY created_at DESC;

-- name: ListVenuesByOwner :many
SELECT id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at
FROM venues
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: GetVenue :one
SELECT id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at
FROM venues
WHERE id = $1;

-- name: CreateVenue :one
INSERT INTO venues (id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
RETURNING id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at;

-- name: UpdateVenue :one
UPDATE venues
SET name = $2,
    description = $3,
    address = $4,
    city = $5,
    capacity = $6,
    contact_phone = $7,
    contact_email = $8,
    status = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, address, city, capacity, contact_phone, contact_email, owner_id, status, created_at, updated_at;

-- name: DeleteVenue :exec
DELETE FROM venues WHERE id = $1;
