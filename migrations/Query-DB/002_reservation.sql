-- Reservation projection table for Query DB
-- This table is populated by the projection worker from events

CREATE TABLE IF NOT EXISTS reservation_projection (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seat_id UUID NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    status TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_reservation_projection_user_id ON reservation_projection(user_id);
CREATE INDEX IF NOT EXISTS idx_reservation_projection_status ON reservation_projection(status);
