package internal

import (
	"bytes"
	"context"
	"html/template"
	"minesweeper/internal/db"
	"minesweeper/internal/models"
)

func GenerateGridHTML(templates *template.Template, game *models.Game) (string, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "game_grid", game)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GetTotalGamesCount(queries *db.Queries) (int64, error) {
	count, err := queries.GetTotalGamesCount(context.Background())

	if err != nil {
		return 0, err
	}

	return count, nil
}
