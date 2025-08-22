CREATE TABLE users (
    id         UUID        PRIMARY KEY,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    name       TEXT        NOT NULL
);
