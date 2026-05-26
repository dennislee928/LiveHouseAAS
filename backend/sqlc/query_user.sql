-- name: GetUserByID :one
SELECT id, email, password_hash, name, role, COALESCE(avatar_url, '') AS avatar_url, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, name, role, COALESCE(avatar_url, '') AS avatar_url, created_at, updated_at
FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, name, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING id, email, name, role, COALESCE(avatar_url, '') AS avatar_url, created_at, updated_at;

-- name: UpdateUser :exec
UPDATE users
SET name = COALESCE($2, name),
    avatar_url = COALESCE($3, avatar_url),
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
