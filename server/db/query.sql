-- name: CreateUser :one
INSERT INTO users (name, password_hash)
VALUES (?, ?)
RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = ? LIMIT 1;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: CreateSession :one
INSERT INTO sessions (
  id, user_id, expires_at
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = ? LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;
