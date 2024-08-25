-- name: GetGames :many 
SELECT * FROM games;

-- name: GetGame :one
SELECT * FROM games WHERE id = ?;

-- name: DeleteGame :exec
DELETE FROM games WHERE id = ?;