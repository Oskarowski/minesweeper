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

-- name: UpdateGameGridStateById :exec
UPDATE
    games
SET
    game_failed = ?,
    game_won = ?,
    grid_state = ?
WHERE
    id = ?;

-- name: GetMovesByGameId :many
SELECT
    *
FROM
    moves
WHERE
    game_id = ?