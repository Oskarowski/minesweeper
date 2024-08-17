package src

import (
	"fmt"
	"strconv"
)

func GenerateGridHTML(game *Game) string {
	html := "<div class='grid gap-1' style='grid-template-columns: repeat(" + strconv.Itoa(game.GridSize) + ", 1fr);'>"

	for row := 0; row < game.GridSize; row++ {
		for col := 0; col < game.GridSize; col++ {
			cell := game.Grid[row][col]

			cellID := fmt.Sprintf("cell-%d-%d", row, col)
			content := ""

			if cell.IsRevealed {
				if cell.HasMine {
					content = "💣"
				} else if cell.AdjacentMines > 0 {
					content = strconv.Itoa(cell.AdjacentMines)
				}
			}

			classes := "border border-gray-400 text-center mine-field aspect-square flex items-center justify-center"

			html += fmt.Sprintf("<div id='%s' class='%s' hx-get='/reveal?row=%d&col=%d'>%s</div>", cellID, classes, row, col, content)

		}

	}

	html += "</div>"

	return html
}

// func RevealCellHandler(w http.ResponseWriter, r *http.Request) {
// 	row, _ := strconv.Atoi(r.URL.Query().Get("row"))
// 	col, _ := strconv.Atoi(r.URL.Query().Get("col"))

// 	//TODO add a way to retrieve the game from the request

// 	mineHit, adjacentMines := currentGame.RevealCell(row, col)

// 	if mineHit {
// 		http.Error(w, "Game Ove! You hit a mine!", http.StatusOK)
// 		return
// 	}

// 	gameGridHtml := gameGenerateGridHTML(currentGame)

// 	w.Header().Set("HX-Trigger", "grid-updated")
// 	w.Header().Set("HX-Trigger-After-Swap", "reveal-complete")
// 	w.Write([]byte(gameGridHtml))
// }
