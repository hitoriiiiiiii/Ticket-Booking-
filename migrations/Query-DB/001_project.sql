
CREATE TABLE IF NOT EXISTS projection_state (
    id SERIAL PRIMARY KEY,
    last_event_id UUID
);

INSERT INTO projection_state (last_event_id)
VALUES (NULL);
