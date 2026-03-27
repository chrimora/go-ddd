CREATE TABLE posts (
    id           UUID        PRIMARY KEY,
    version      INT         NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL,
    title        TEXT        NOT NULL,
    publish_date TIMESTAMPTZ NOT NULL,
    author       TEXT        NOT NULL
);

CREATE INDEX ON posts (id);
