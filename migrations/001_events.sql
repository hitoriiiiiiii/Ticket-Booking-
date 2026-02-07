CREATE EXTENSION IF NOT EXISTS "pgcrypto";

DROP TABLE IF EXISTS events;

CREATE TABLE events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  aggregate_id UUID NOT NULL,
  event_type TEXT NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);
