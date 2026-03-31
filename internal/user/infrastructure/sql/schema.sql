CREATE TABLE users (
    id         UUID        PRIMARY KEY,
    version    INT         NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name       TEXT        NOT NULL
);

CREATE INDEX ON users (id);
