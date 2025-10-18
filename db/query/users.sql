-- name: CreateUser :one
INSERT INTO users (name, email, phone, balance, created_at)
VALUES ($1, $2, $3, 0.00, NOW())
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET name=$2, email=$3, phone=$4
WHERE id=$1
RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name=$1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email=$1
LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id=$1;
