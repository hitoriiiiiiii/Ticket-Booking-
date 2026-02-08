CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  aggregate_id TEXT NOT NULL,
  event_type TEXT NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE reservations (
    seat_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL,
    reserved_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE projection_state (
    id SERIAL PRIMARY KEY,
    last_event_id UUID
);

INSERT INTO projection_state (last_event_id) VALUES (NULL);
