package internal

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"math/rand"
	"minesweeper/internal/db"
	"minesweeper/internal/models"
	"strconv"
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

type GameSettings struct {
	GridSize    int
	MinesAmount int
}

func ValidateGameSettingsForm(gridSizeStr, minesAmountStr, randomMinesStr, randomGridSizeStr string) (GameSettings, error) {

	gridSize, gridSizeErr := strconv.Atoi(gridSizeStr)
	minesAmount, minesAmountErr := strconv.Atoi(minesAmountStr)

	if randomGridSizeStr == "on" {
		// TODO - load min and max grid size from config
		gridSize = rand.Intn(10) + 5 // range(5, 15)
		gridSizeErr = nil
	}

	if randomMinesStr == "on" {
		// TODO - mines amount based on chosen difficulty?
		minesAmount = rand.Intn((gridSize*gridSize)/2) + 1
		minesAmountErr = nil
	}

	if gridSizeErr != nil || minesAmountErr != nil {
		return GameSettings{}, errors.New("invalid input values")
	}

	if gridSize <= 0 || minesAmount <= 0 {
		return GameSettings{}, errors.New("the grid size and mines amount must be greater than 0")
	}

	return GameSettings{
		GridSize:    gridSize,
		MinesAmount: minesAmount,
	}, nil
}
