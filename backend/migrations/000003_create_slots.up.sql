CREATE TABLE slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id UUID NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('available','booked','blocked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_slots_venue ON slots(venue_id);
CREATE INDEX idx_slots_date ON slots(date);
CREATE INDEX idx_slots_status ON slots(status);
CREATE UNIQUE INDEX idx_slots_venue_time ON slots(venue_id, date, start_time, end_time);
