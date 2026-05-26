CREATE TABLE seat_layouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id UUID NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL DEFAULT 'Main',
    rows INT NOT NULL DEFAULT 10,
    cols INT NOT NULL DEFAULT 10,
    seats JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_seat_layout_venue ON seat_layouts(venue_id);

ALTER TABLE ticket_types ADD COLUMN IF NOT EXISTS seat_section VARCHAR(100) DEFAULT '';
ALTER TABLE ticket_types ADD COLUMN IF NOT EXISTS seat_rows VARCHAR(50) DEFAULT '';
