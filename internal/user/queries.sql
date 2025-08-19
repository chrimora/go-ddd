-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (name)
VALUES ($1)
RETURNING id;

-- name: UpdateUser :exec
UPDATE users
SET name = $3, updated_at = now()
WHERE id = $1
AND updated_at = $2
RETURNING id;
