-- name: ClaimNextEventBatch :many
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
    AND o1.status = 'Pending'
    AND o1.retries < $2
    ORDER BY o1.created_at, o1.aggregate_id
    FOR UPDATE SKIP LOCKED
    LIMIT $1
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

-- name: RequeueStaleEvents :many
UPDATE event_outbox
SET status = 'Pending', retries = retries + 1, updated_at = now()
WHERE status = 'Claimed'
AND updated_at < $1
AND retries < $2
RETURNING id;

-- name: CompleteEvent :exec
UPDATE event_outbox
SET status = 'Published', updated_at = now()
WHERE id = $1;
