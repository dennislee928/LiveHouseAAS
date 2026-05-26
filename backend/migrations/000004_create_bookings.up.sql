CREATE TABLE booking_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id UUID NOT NULL REFERENCES slots(id),
    venue_id UUID NOT NULL REFERENCES venues(id),
    artist_id UUID NOT NULL REFERENCES users(id),
    message TEXT DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','approved','rejected','cancelled','confirmed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bookings_venue ON booking_requests(venue_id);
CREATE INDEX idx_bookings_artist ON booking_requests(artist_id);
CREATE INDEX idx_bookings_status ON booking_requests(status);

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    venue_id UUID NOT NULL REFERENCES venues(id),
    artist_id UUID NOT NULL REFERENCES users(id),
    booking_id UUID REFERENCES booking_requests(id),
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','published','cancelled','completed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_venue ON events(venue_id);
CREATE INDEX idx_events_artist ON events(artist_id);
CREATE INDEX idx_events_status ON events(status);
