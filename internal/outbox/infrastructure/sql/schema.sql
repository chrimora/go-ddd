CREATE TYPE event_status AS ENUM ('Pending', 'Claimed', 'Published');

CREATE TABLE event_outbox (
    id             UUID         PRIMARY KEY,
    aggregate_id   UUID         NOT NULL,
    aggregate_type VARCHAR(32)  NOT NULL,
    event_context  JSONB        NOT NULL,
    event_type     VARCHAR(64)  NOT NULL,
    payload        JSONB        NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL,
    updated_at     TIMESTAMPTZ  NOT NULL,
    retries        INT          NOT NULL,
    status         event_status NOT NULL
);

CREATE INDEX idx_event_outbox_aggregate_id ON event_outbox (aggregate_id);
