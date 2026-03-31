CREATE TABLE orders (
    id         UUID        PRIMARY KEY,
    version    INT         NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status     TEXT        NOT NULL
);

CREATE INDEX ON orders (id);

CREATE TABLE order_items (
    id         UUID   PRIMARY KEY,
    order_id   UUID   NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    name       TEXT   NOT NULL,
    quantity   INT    NOT NULL,
    unit_price BIGINT NOT NULL
);

CREATE INDEX ON order_items (order_id);
