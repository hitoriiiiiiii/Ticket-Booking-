
CREATE TABLE IF NOT EXISTS reservations (
    seat_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL,
    reserved_at TIMESTAMP DEFAULT NOW()
);