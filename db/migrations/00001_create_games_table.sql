-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    games (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT NOT NULL DEFAULT (lower(hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-4' || substr(hex(randomblob(2)),2) || '-' || substr('AB89',abs(random()) % 4 + 1, 1) || substr(hex(randomblob(2)),2) || '-' || hex(randomblob(6)))),
        grid_size INTEGER NOT NULL,
        mines_amount INTEGER NOT NULL,
        game_failed BOOLEAN NOT NULL DEFAULT FALSE,
        game_won BOOLEAN NOT NULL DEFAULT FALSE,
        grid_state TEXT NOT NULL,
        create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE games
-- +goose StatementEnd
-- goose -dir db/migrations sqlite3 ./db/minesweeper.db up