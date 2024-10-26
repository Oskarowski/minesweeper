package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

const (
	MinGridSize   = 2
	MaxGridSize   = 22
	MinMinesRatio = 0.1
	MaxMinesRatio = 0.8
)

type GameSettings struct {
	GridSize    int
	MinesAmount int
}

func ValidateGameSettingsForm(gridSizeStr, minesAmountStr, randomMinesStr, randomGridSizeStr string) (GameSettings, error) {
	var (
		gridSize, minesAmount       int
		gridSizeErr, minesAmountErr error
	)

	// Check if grid size should be random or user-defined, if so check if it's within accepted bounds
	if randomGridSizeStr == "on" {
		gridSize = rand.Intn(MaxGridSize-MinGridSize+1) + MinGridSize
		gridSizeErr = nil
	} else {
		gridSize, gridSizeErr = strconv.Atoi(gridSizeStr)
		if gridSizeErr != nil {
			return GameSettings{}, errors.New("invalid grid size: must be a proper grid size number")
		}

		if gridSize < MinGridSize || gridSize > MaxGridSize {
			return GameSettings{}, errors.New("grid size must be between 2 and 50")
		}

	}

	minMines := int(float64(gridSize*gridSize) * MinMinesRatio)
	maxMines := int(float64(gridSize*gridSize) * MaxMinesRatio)

	// Check if mines amount should be random or user-defined, if so check if it's within accepted bounds
	if randomMinesStr == "on" {
		if gridSize > 0 {

			minesAmount = rand.Intn((maxMines - minMines + 1)) + minMines
			minesAmountErr = nil
		} else {
			return GameSettings{}, errors.New("grid size must be valid when using random mines")
		}
	} else {
		minesAmount, minesAmountErr = strconv.Atoi(minesAmountStr)

		if minesAmountErr != nil {
			return GameSettings{}, errors.New("invalid mines amount: must be a number")
		}

		if minesAmount <= 0 || minesAmount > maxMines {
			return GameSettings{}, fmt.Errorf("mines amount must be between 1 and %v of the grid size", maxMines)
		}
	}

	if gridSizeErr != nil || minesAmountErr != nil {
		return GameSettings{}, errors.New("invalid input values")
	}

	return GameSettings{
		GridSize:    gridSize,
		MinesAmount: minesAmount,
	}, nil
}
