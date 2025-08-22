-- name: GetNextEvent :one
SELECT *
FROM event_outbox
WHERE processed_at IS NULL
AND retries < 3
ORDER BY created_at
LIMIT 1
FOR UPDATE SKIP LOCKED;

-- name: CreateEvent :exec
INSERT INTO event_outbox (event_type, payload)
VALUES ($1, $2);

-- name: UpdateEvent :one
UPDATE event_outbox
SET updated_at = now(),
    retries = $2,
    processed_at = $3
WHERE id = $1
RETURNING id;
