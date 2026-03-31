-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1;

-- name: CreateOrder :one
INSERT INTO orders (id, version, status)
VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateOrder :one
UPDATE orders
SET version = version + 1, updated_at = NOW(), status = $3
WHERE id = $1
AND version = $2
RETURNING id;

-- name: RemoveOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: GetOrderItems :many
SELECT * FROM order_items WHERE order_id = $1;

-- name: CreateOrderItem :exec
INSERT INTO order_items (id, order_id, name, quantity, unit_price)
VALUES ($1, $2, $3, $4, $5);

-- name: DeleteOrderItems :exec
DELETE FROM order_items WHERE order_id = $1;
