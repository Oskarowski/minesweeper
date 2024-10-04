package models

import (
	"log"
	"math/rand"
	"minesweeper/internal/db"
	"strings"
)

type Cell struct {
	HasMine       bool
	IsFlagged     bool
	IsRevealed    bool
	AdjacentMines int
}

type Game struct {
	Id          int64
	Uuid        string
	GridSize    int
	MinesAmount int
	Grid        [][]Cell
	GameFailed  bool
	GameWon     bool
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

func (g *Game) RevealCell(row int, col int) {
	if row < 0 || row >= g.GridSize || col < 0 || col >= g.GridSize {
		return
	}

	cell := &g.Grid[row][col]

	if cell.IsFlagged {
		return
	}

	if cell.IsRevealed {
		return
	}

	cell.IsRevealed = true

	if cell.HasMine {
		log.Printf("Game with UUID: %s is ended and FAILED", g.Uuid)

		g.GameFailed = true
		return
	}
	g.CheckWinCondition()

	if cell.AdjacentMines > 0 {
		return
	}

	g.revealSurroundingCells(row, col)
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

	g.CheckWinCondition()
}

func (g *Game) CheckWinCondition() bool {
	if g.GameWon {
		return g.GameWon
	}

	var flaggedMines uint16 = 0
	var revealedCells uint16 = 0
	totalMines := uint16(g.MinesAmount)

	for _, row := range g.Grid {
		for _, cell := range row {
			if cell.HasMine && cell.IsFlagged {
				flaggedMines++
			}

			if cell.IsRevealed {
				revealedCells++
			}

			if !cell.HasMine && cell.IsFlagged {
				return false
			}
		}
	}

	var totalCells = uint16(g.GridSize * g.GridSize)
	var allEmptyCellsRevealed bool = totalCells-flaggedMines == revealedCells
	var flaggedAllMines = totalMines == flaggedMines

	if allEmptyCellsRevealed && flaggedAllMines {
		log.Printf("Game with UUID: %s is ended and WON", g.Uuid)
		g.GameWon = true
		return true
	}

	return false
}

const (
	CELL_REVEALED     = 'R'
	CELL_FLAGGED      = 'F'
	CELL_FLAGGED_MINE = 'X'
	CELL_HAS_MINE     = 'M'
	CELL_EMPTY        = 'E'
	ROW_SEPARATOR     = '|'
)

// EncodeGameGrid encodes a game grid into a string for easier storage and transport.
// The encoding is not taking into account the adjacent mines of each cell because it can be easily computed.
// The encoding is as follows:
//
// - 'R' for a revealed cell
// - 'F' for a flagged cell
// - 'X' for a flagged cell with a mine
// - 'M' for a non-revealed cell with a mine
// - 'E' for a non-revealed cell without a mine
//
// The rows are separated by '|' characters.
func EncodeGameGrid(grid [][]Cell) string {
	var sb strings.Builder

	for _, row := range grid {
		for _, cell := range row {
			if cell.IsRevealed {
				sb.WriteRune(CELL_REVEALED)
			} else if cell.IsFlagged && cell.HasMine {
				sb.WriteRune(CELL_FLAGGED_MINE)
			} else if cell.HasMine {
				sb.WriteRune(CELL_HAS_MINE)
			} else if cell.IsFlagged {
				sb.WriteRune(CELL_FLAGGED)
			} else {
				sb.WriteRune(CELL_EMPTY)
			}
		}
		// the `|` is used as a delimiter between rows
		sb.WriteRune(ROW_SEPARATOR)
	}

	return sb.String()
}

// DecodeGameGrid decodes a game grid from a string.
//
// The decoding is the reverse of EncodeGameGrid, i.e. it takes the string
// representation of the game grid and returns a 2D slice of Cell structs.
//
// Also calculates the number of adjacent mines for each cell in the grid with call to updateAdjacentMines.
func DecodeGameGrid(encodedGameGrid string, gridSize int) [][]Cell {
	rows := strings.Split(encodedGameGrid, "|")
	decodedGameGrid := make([][]Cell, gridSize)

	if len(encodedGameGrid) == 0 {
		panic("encodedGameGrid == 0")
	}

	for i, row := range rows {
		if len(row) == 0 {
			break
		}

		decodedGameGrid[i] = make([]Cell, gridSize)

		for j, char := range row {
			switch char {
			case CELL_REVEALED:
				decodedGameGrid[i][j] = Cell{IsRevealed: true}
			case CELL_FLAGGED_MINE:
				decodedGameGrid[i][j] = Cell{IsFlagged: true, HasMine: true}
			case CELL_FLAGGED:
				decodedGameGrid[i][j] = Cell{IsFlagged: true}
			case CELL_HAS_MINE:
				decodedGameGrid[i][j] = Cell{HasMine: true}
			case CELL_EMPTY:
				decodedGameGrid[i][j] = Cell{AdjacentMines: 0}
			}
		}
	}

	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize; col++ {
			if decodedGameGrid[row][col].HasMine {
				updateAdjacentCells(decodedGameGrid, row, col, gridSize)
			}
		}
	}

	return decodedGameGrid
}

func FromDbGame(dbGame *db.Game) (*Game, error) {
	decodedGameGrid := DecodeGameGrid(dbGame.GridState, int(dbGame.GridSize))

	return &Game{
		Id:          dbGame.Id,
		Uuid:        dbGame.Uuid,
		GridSize:    int(dbGame.GridSize),
		MinesAmount: int(dbGame.MinesAmount),
		Grid:        decodedGameGrid,
		GameFailed:  dbGame.GameFailed,
		GameWon:     dbGame.GameWon,
	}, nil
}

func ToDbGame(game *Game) (*db.Game, error) {
	// encodedGameGrid := EncodeGameGrid(game.Grid)

	// return &db.Game{
	// 	ID:          game.ID,
	// 	Uuid:        game.Uuid,
	// 	GridSize:    int64(game.GridSize),
	// 	MinesAmount: int64(game.MinesAmount),
	// 	GridState:   encodedGameGrid,
	// 	GameFailed:  game.GameFailed,
	// 	GameWon:     game.GameWon,
	// 	CreateAt:    sql.NullTime{Time: time.Now(), Valid: true},
	// }, nil
	return nil, nil
}
