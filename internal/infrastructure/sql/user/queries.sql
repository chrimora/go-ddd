-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, version, created_at, updated_at, name)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateUser :one
UPDATE users
SET version = version + 1, updated_at = $3, name = $4 
WHERE id = $1
AND version = $2
RETURNING id;
