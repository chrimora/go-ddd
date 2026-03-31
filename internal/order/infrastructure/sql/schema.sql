CREATE TABLE orders (
    id         UUID        PRIMARY KEY,
    version    INT         NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status     TEXT        NOT NULL,
    total      BIGINT      NOT NULL
);

CREATE INDEX ON orders (id);
