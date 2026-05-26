CREATE TABLE nft_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    token_id BIGINT NOT NULL,
    contract_address VARCHAR(255) NOT NULL DEFAULT '',
    tx_hash VARCHAR(255) DEFAULT '',
    network VARCHAR(50) NOT NULL DEFAULT 'polygon',
    token_uri TEXT DEFAULT '',
    owner_address VARCHAR(255) NOT NULL DEFAULT '',
    is_poap BOOLEAN NOT NULL DEFAULT false,
    poap_claimed_at TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','minted','failed','claimed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_nft_ticket ON nft_tickets(ticket_id);
CREATE INDEX idx_nft_owner ON nft_tickets(owner_address);
