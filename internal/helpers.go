package internal

import (
	"bytes"
	"html/template"
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
