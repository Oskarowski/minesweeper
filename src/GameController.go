package src

import (
	"math/rand"
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
	GameFailed  bool
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
		for j := range grid[i] {
			grid[i][j] = Cell{
				IsRevealed:    false,
				HasMine:       false,
				AdjacentMines: 0,
			}
		}
	}

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

func (g *Game) revealSurroundingCells(row int, col int) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue // we are skipping the current cell
			}

			r, c := row+i, col+j

			if r >= 0 && r < g.GridSize && c >= 0 && c < g.GridSize {
				neighbor := &g.Grid[r][c]

				if !neighbor.IsRevealed && !neighbor.HasMine {
					neighbor.IsRevealed = true
					if neighbor.AdjacentMines == 0 {
						g.revealSurroundingCells(r, c)
					}
				}
			}
		}
	}
}

func (g *Game) RevealCell(row int, col int) (bool, int) {
	if row < 0 || row >= g.GridSize || col < 0 || col >= g.GridSize {
		return false, 0
	}

	cell := &g.Grid[row][col]

	if cell.IsFlagged {
		return false, 0
	}

	if cell.IsRevealed {
		return cell.HasMine, cell.AdjacentMines
	}

	cell.IsRevealed = true

	if cell.HasMine {
		g.GameFailed = true
		return true, -1
	}

	if cell.AdjacentMines > 0 {
		return false, cell.AdjacentMines
	}

	g.revealSurroundingCells(row, col)

	return false, 0
}

func (g *Game) FlagCell(row int, col int) {
	if row < 0 || row >= g.GridSize || col < 0 || col >= g.GridSize {
		return
	}
	cell := &g.Grid[row][col]

	if cell.IsRevealed {
		return
	}

	cell.IsFlagged = !cell.IsFlagged
}
