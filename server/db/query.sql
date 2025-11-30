-- name: GetGameStats :many
SELECT * FROM game_stats
ORDER BY created_at DESC;

-- name: CreateGameStat :one
INSERT INTO game_stats (created_at)
VALUES (CURRENT_TIMESTAMP)
RETURNING *;
