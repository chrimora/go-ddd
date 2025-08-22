CREATE TABLE event_outbox (
    id           SERIAL      PRIMARY KEY,
    event_type   VARCHAR(64) NOT NULL,
    payload      JSONB       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    retries      INT         NOT NULL DEFAULT 0,
    processed_at TIMESTAMPTZ
);

