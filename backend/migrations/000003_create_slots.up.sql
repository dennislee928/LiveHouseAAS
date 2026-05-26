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

-- prevent overlapping slots for the same venue
CREATE EXTENSION IF NOT EXISTS btree_gist;
CREATE CONSTRAINT EXCLUDE USING gist (
    venue_id WITH =,
    daterange(date, date, '[]') WITH &&
    tsrange(
        date + start_time,
        date + end_time,
        '[]'
    ) WITH &&
);
