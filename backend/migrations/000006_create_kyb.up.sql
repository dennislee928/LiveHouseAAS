CREATE TABLE business_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    business_name VARCHAR(255) NOT NULL,
    tax_id VARCHAR(50) NOT NULL,
    registration_number VARCHAR(100) DEFAULT '',
    address TEXT NOT NULL,
    phone VARCHAR(50) DEFAULT '',
    documents JSONB DEFAULT '[]',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','verified','rejected')),
    rejection_reason TEXT DEFAULT '',
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_kyb_user ON business_verifications(user_id);
CREATE INDEX idx_kyb_status ON business_verifications(status);

ALTER TABLE venues ADD COLUMN IF NOT EXISTS verified_only BOOLEAN NOT NULL DEFAULT false;
