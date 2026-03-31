-- name: ListOrders :many
SELECT *
FROM orders
WHERE (sqlc.narg(after)::uuid IS NULL OR id > sqlc.narg(after))
  AND (sqlc.narg(status)::text IS NULL OR status = sqlc.narg(status))
ORDER BY id
LIMIT @limit_plus_one;
