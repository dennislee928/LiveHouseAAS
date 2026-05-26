-- name: ListVenueSpecs :many
SELECT id, venue_id, category, name, brand, quantity, description, created_at, updated_at
FROM venue_specs
WHERE venue_id = $1
ORDER BY category, name;

-- name: GetVenueSpec :one
SELECT id, venue_id, category, name, brand, quantity, description, created_at, updated_at
FROM venue_specs
WHERE id = $1;

-- name: CreateVenueSpec :one
INSERT INTO venue_specs (id, venue_id, category, name, brand, quantity, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
RETURNING id, venue_id, category, name, brand, quantity, description, created_at, updated_at;

-- name: UpdateVenueSpec :one
UPDATE venue_specs
SET category = $2,
    name = $3,
    brand = $4,
    quantity = $5,
    description = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING id, venue_id, category, name, brand, quantity, description, created_at, updated_at;

-- name: DeleteVenueSpec :exec
DELETE FROM venue_specs WHERE id = $1;

-- name: DeleteVenueSpecsByVenue :exec
DELETE FROM venue_specs WHERE venue_id = $1;
