CREATE TABLE event_outbox (
    id            SERIAL      PRIMARY KEY,
    aggregate_id  UUID        NOT NULL,
    event_context JSONB       NOT NULL,
    event_type    VARCHAR(64) NOT NULL,
    payload       JSONB       NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    retries       INT         NOT NULL DEFAULT 0,
    processed_at  TIMESTAMPTZ
);

CREATE INDEX idx_event_outbox_aggregate_id ON event_outbox (aggregate_id);
