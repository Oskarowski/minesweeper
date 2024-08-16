package src

import (
	"math/rand"
	"time"
)

type Cell struct {
	HasMine       bool
	IsFlagged     bool
	IsRevealed    bool
	AdjacentMines int
}

type Game struct {
	GridSize    int
	MinesAmount int
	Grid        [][]Cell
}

func updateAdjacentCells(grid [][]Cell, row int, col int, gridSize int) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			r, c := row+i, col+j
			if r >= 0 && r < gridSize && c >= 0 && c < gridSize {
				grid[r][c].AdjacentMines++
			}
		}
	}
}

func NewGame(gridSize int, minesAmount int) *Game {
	grid := make([][]Cell, gridSize)

	for i := range grid {
		grid[i] = make([]Cell, gridSize)
	}

	rand.Seed(time.Now().UnixNano())
	minesPlaced := 0

	for minesPlaced < minesAmount {
		row := rand.Intn(gridSize)
		col := rand.Intn(gridSize)

		if !grid[row][col].HasMine {
			grid[row][col].HasMine = true
			minesPlaced++
			updateAdjacentCells(grid, row, col, gridSize)
		}
	}

	return &Game{
		GridSize:    gridSize,
		MinesAmount: minesAmount,
		Grid:        grid,
	}
}

func (g *Game) RevealCell(row int, col int) (bool, int) {
	if row < 0 || row >= g.GridSize || col < 0 || col >= g.GridSize {
		return false, 0
	}

	cell := &g.Grid[row][col]

	if cell.IsRevealed {
		return cell.HasMine, cell.AdjacentMines
	}

	cell.IsRevealed = true
	if cell.HasMine {
		return true, -1
	}

	return false, cell.AdjacentMines
}
