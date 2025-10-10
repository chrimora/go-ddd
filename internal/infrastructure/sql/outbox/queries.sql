-- name: ClaimNextEvent :one
WITH candidates AS (
    SELECT o1.id
    FROM event_outbox o1
    WHERE NOT EXISTS (
        SELECT 1
        FROM event_outbox o2
        WHERE o2.aggregate_id = o1.aggregate_id
        AND o2.created_at < o1.created_at
        AND o2.status IN ('Pending', 'Claimed')
    )
    AND status = 'Pending'
    AND retries < 3
    ORDER BY created_at, aggregate_id
    FOR UPDATE SKIP LOCKED
    LIMIT 1 -- increase to collect a batch
)
UPDATE event_outbox
SET status = 'Claimed', updated_at = now()
WHERE id IN (SELECT id FROM candidates)
RETURNING *;

-- name: CreateEvent :exec
INSERT INTO event_outbox (
    id,
    aggregate_id,
    aggregate_type,
    event_context,
    event_type,
    payload,
    created_at,
    updated_at,
    retries,
    status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;

-- name: RequeueEvent :exec
UPDATE event_outbox
SET status = 'Pending', retries = retries + 1, updated_at = now()
WHERE id = $1;

-- name: CompleteEvent :exec
UPDATE event_outbox
SET status = 'Processed', updated_at = now()
WHERE id = $1;
