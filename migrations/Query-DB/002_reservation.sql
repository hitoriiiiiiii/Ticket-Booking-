
CREATE TABLE IF NOT EXISTS reservation_projection (
    seat_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);