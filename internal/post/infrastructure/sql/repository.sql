-- name: GetPost :one
SELECT * FROM posts WHERE id = $1;

-- name: CreatePost :one
INSERT INTO posts (id, version, created_at, updated_at, title, publish_date, author)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: UpdatePost :one
UPDATE posts
SET version = version + 1, updated_at = $3, title = $4
WHERE id = $1
AND version = $2
RETURNING id;

-- name: RemovePost :exec
DELETE FROM posts
WHERE id = $1;
