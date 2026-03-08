
-- Create new migration for updated reservations table with UUID support
-- This migration updates the reservations table to use UUIDs and adds constraints

-- Drop existing table if exists and recreate with proper UUID columns
DROP TABLE IF EXISTS reservations CASCADE

CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    seat_id UUID NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    show_id UUID,
    status TEXT NOT NULL DEFAULT 'HELD',
    reserved_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add foreign key constraint to users
ALTER TABLE reservations 
ADD CONSTRAINT fk_reservation_user 
FOREIGN KEY (user_id) REFERENCES users(id);

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_reservations_user_id ON reservations(user_id);
CREATE INDEX IF NOT EXISTS idx_reservations_seat_id ON reservations(seat_id);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);
