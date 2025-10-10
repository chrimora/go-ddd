CREATE TABLE users (
    id         UUID        PRIMARY KEY,
    version    INT         NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    name       TEXT        NOT NULL
);
