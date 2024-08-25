-- +goose Up
-- +goose StatementBegin
CREATE TABLE games (
    id INT PRIMARY KEY,
    grid_size INT NOT NULL,
    mines_amount INT NOT NULL,
    game_failed BOOLEAN,
    game_won BOOLEAN,
    create_at TIMESTAMP DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE games
-- +goose StatementEnd

-- goose -dir db/migrations sqlite3 ./db/minesweeper.db up