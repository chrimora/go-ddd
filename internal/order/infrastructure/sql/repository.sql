-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (id, version, status, total)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: UpdateOrder :one
UPDATE orders
SET version = version + 1, updated_at = NOW(), status = $3, total = $4
WHERE id = $1
AND version = $2
RETURNING id;

-- name: RemoveOrder :exec
DELETE FROM orders
WHERE id = $1;
