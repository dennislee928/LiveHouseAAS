-- name: GetNFTByTicket :one
SELECT id, ticket_id, token_id, contract_address, tx_hash, network, token_uri, owner_address, is_poap, poap_claimed_at, status, created_at, updated_at
FROM nft_tickets
WHERE ticket_id = $1;

-- name: CreateNFTRecord :one
INSERT INTO nft_tickets (id, ticket_id, token_id, contract_address, tx_hash, network, token_uri, owner_address, is_poap, poap_claimed_at, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, false, NULL, 'pending', NOW(), NOW())
RETURNING id, ticket_id, token_id, contract_address, tx_hash, network, token_uri, owner_address, is_poap, poap_claimed_at, status, created_at, updated_at;

-- name: UpdateNFTStatus :one
UPDATE nft_tickets
SET status = $2, tx_hash = COALESCE($3, tx_hash), token_uri = COALESCE($4, token_uri), updated_at = NOW()
WHERE id = $1
RETURNING id, ticket_id, token_id, contract_address, tx_hash, network, token_uri, owner_address, is_poap, poap_claimed_at, status, created_at, updated_at;

-- name: ClaimPOAP :one
UPDATE nft_tickets
SET is_poap = true, poap_claimed_at = NOW(), status = 'claimed', updated_at = NOW()
WHERE ticket_id = $1 AND status = 'minted'
RETURNING id, ticket_id, token_id, contract_address, tx_hash, network, token_uri, owner_address, is_poap, poap_claimed_at, status, created_at, updated_at;
