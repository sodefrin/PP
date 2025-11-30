-- name: CreateUser :one
INSERT INTO users (name, password_hash)
VALUES (?, ?)
RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = ? LIMIT 1;
