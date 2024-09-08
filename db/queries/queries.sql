-- name: CreateGame :one 
INSERT INTO
    games (grid_size, mines_amount, grid_state)
VALUES
    (?, ?, ?) RETURNING *;

-- name: InsertMove :one
INSERT INTO
    moves (game_id, move_type, row, col)
VALUES
    (?, ?, ?, ?) RETURNING *;

-- name: GetGameById :one
SELECT
    *
FROM
    games
WHERE
    id = ?;

-- name: GetGameByUuid :one
SELECT
    *
FROM
    games
WHERE
    uuid = ?;

-- name: ListGames :many
SELECT
    id, grid_size, mines_amount, game_failed, game_won, created_at
FROM
    games
ORDER BY
    created_at DESC
LIMIT ? OFFSET ?;

-- name: GetTotalGamesCount :one
SELECT
    COUNT(id) as count
FROM
    games;

-- name: UpdateGameGridStateById :exec
UPDATE
    games
SET
    game_failed = ?,
    game_won = ?,
    grid_state = ?
WHERE
    id = ?;

-- name: GetGamesInfoByUuids :one
SELECT 
    COUNT(*) AS total_games,
    COUNT(*) FILTER (WHERE game_won = TRUE) AS won_games,
    COUNT(*) FILTER (WHERE game_failed = TRUE AND game_won = FALSE) AS lost_games,
    COUNT(*) FILTER (WHERE game_failed = FALSE AND game_won = FALSE) AS not_finished_games
FROM 
    games
WHERE 
    uuid IN (sqlc.slice('uuids'));

-- name: GetMovesByGameId :many
SELECT
    *
FROM
    moves
WHERE
    game_id = ?