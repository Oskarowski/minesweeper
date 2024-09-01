-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    moves (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        game_id INTEGER NOT NULL,
        move_type VARCHAR(255) NOT NULL,
        row INTEGER NOT NULL,
        col INTEGER NOT NULL,
        create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (game_id) REFERENCES games (id)
    )
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE moves
-- +goose StatementEnd