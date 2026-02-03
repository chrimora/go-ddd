-- name: ListUsers :many
SELECT *
FROM users
WHERE (sqlc.narg(after)::uuid IS NULL OR id > sqlc.narg(after))
  AND (sqlc.narg(name)::text IS NULL OR name ILIKE '%' || sqlc.narg(name) || '%')
ORDER BY id
LIMIT @limit_plus_one;
