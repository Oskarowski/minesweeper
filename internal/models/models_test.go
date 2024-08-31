package models

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEncodeGameGrid(t *testing.T) {
	grid := [][]Cell{
		{{HasMine: true}, {IsRevealed: true}, {IsFlagged: true}, {AdjacentMines: 1}},
		{{}, {HasMine: true}, {}, {IsFlagged: true}},
	}

	expected := string(CELL_HAS_MINE) + string(CELL_REVELED) + string(CELL_FLAGGED) + string(CELL_EMPTY) + string(ROW_SEPARATOR) + string(CELL_EMPTY) + string(CELL_HAS_MINE) + string(CELL_EMPTY) + string(CELL_FLAGGED) + string(ROW_SEPARATOR)

	var encodedGrid string = EncodeGameGrid(grid)

	if encodedGrid != expected {
		t.Errorf("Expected encoded grid to be '%s', but got '%s'", expected, encodedGrid)
	}
}

func TestDecodeGameGrid(t *testing.T) {
	encodedGameString := "MR|EM|"

	expected := [][]Cell{
		{{HasMine: true, AdjacentMines: 1}, {IsRevealed: true, AdjacentMines: 2}},
		{{AdjacentMines: 2}, {HasMine: true, AdjacentMines: 1}},
	}
	fmt.Println(encodedGameString)

	decoded := DecodeGameGrid(encodedGameString, 2)

	if !reflect.DeepEqual(decoded, expected) {
		t.Errorf("Expected decoded grid to be '%v', but got '%v'", expected, decoded)
	}
}
