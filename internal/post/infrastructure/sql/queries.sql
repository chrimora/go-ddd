-- name: ListPosts :many
SELECT *
FROM posts
WHERE (sqlc.narg(after)::uuid IS NULL OR id > sqlc.narg(after))
  AND (sqlc.narg(title)::text IS NULL OR title ILIKE '%' || sqlc.narg(title) || '%')
ORDER BY id
LIMIT @limit_plus_one;
