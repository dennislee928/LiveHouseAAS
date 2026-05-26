-- name: GetBusinessVerification :one
SELECT id, user_id, business_name, tax_id, registration_number, address, phone, documents, status, COALESCE(rejection_reason,''), verified_at, created_at, updated_at
FROM business_verifications
WHERE user_id = $1;

-- name: CreateBusinessVerification :one
INSERT INTO business_verifications (id, user_id, business_name, tax_id, registration_number, address, phone, documents, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending', NOW(), NOW())
RETURNING id, user_id, business_name, tax_id, registration_number, address, phone, documents, status, COALESCE(rejection_reason,''), verified_at, created_at, updated_at;

-- name: UpdateBusinessVerificationStatus :one
UPDATE business_verifications
SET status = $2, rejection_reason = $3, verified_at = CASE WHEN $2 = 'verified' THEN NOW() ELSE verified_at END, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, business_name, tax_id, registration_number, address, phone, documents, status, COALESCE(rejection_reason,''), verified_at, created_at, updated_at;

-- name: ListPendingVerifications :many
SELECT id, user_id, business_name, tax_id, registration_number, address, phone, documents, status, COALESCE(rejection_reason,''), verified_at, created_at, updated_at
FROM business_verifications
WHERE status = 'pending'
ORDER BY created_at;
