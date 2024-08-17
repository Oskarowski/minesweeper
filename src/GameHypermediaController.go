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

			html += fmt.Sprintf("<div id='%s' class='%s' hx-get='/reveal?row=%d&col=%d' hx-target='#%s' hx-swap='innerHTML'>%s</div>", cellID, classes, row, col, cellID, content)

		}

	}

	html += "</div>"

	return html
}

func GetCellContent(cell *Cell) string {
	if cell.HasMine {
		return "💣"
	} else if cell.AdjacentMines > 0 {
		return strconv.Itoa(cell.AdjacentMines)
	}
	return ""
}
