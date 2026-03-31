-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, version, name)
VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateUser :one
UPDATE users
SET version = version + 1, updated_at = NOW(), name = $3
WHERE id = $1
AND version = $2
RETURNING id;

-- name: RemoveUser :exec
DELETE FROM users
WHERE id = $1;
