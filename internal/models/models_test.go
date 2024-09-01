package models

import (
	"reflect"
	"testing"
)

func TestEncodeGameGrid(t *testing.T) {
	testCases := []struct {
		name     string
		grid     [][]Cell
		expected string
	}{
		{
			name: "Simple grid with mine, revealed, flagged, empty",
			grid: [][]Cell{
				{{HasMine: true}, {IsRevealed: true}, {IsFlagged: true}, {AdjacentMines: 1}},
				{{}, {HasMine: true}, {}, {IsFlagged: true}},
			},
			expected: string(CELL_HAS_MINE) + string(CELL_REVEALED) + string(CELL_FLAGGED) + string(CELL_EMPTY) + string(ROW_SEPARATOR) +
				string(CELL_EMPTY) + string(CELL_HAS_MINE) + string(CELL_EMPTY) + string(CELL_FLAGGED) + string(ROW_SEPARATOR),
		},
		{
			name: "All cells flagged with mines",
			grid: [][]Cell{
				{{IsFlagged: true, HasMine: true}, {IsFlagged: true, HasMine: true}},
				{{IsFlagged: true, HasMine: true}, {IsFlagged: true, HasMine: true}},
			},
			expected: string(CELL_FLAGGED_MINE) + string(CELL_FLAGGED_MINE) + string(ROW_SEPARATOR) +
				string(CELL_FLAGGED_MINE) + string(CELL_FLAGGED_MINE) + string(ROW_SEPARATOR),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encodedGameGrid := EncodeGameGrid(tc.grid)
			if encodedGameGrid != tc.expected {
				t.Errorf("Test case '%s' failed. Expected encoded grid to be '%s', but got '%s'", tc.name, tc.expected, encodedGameGrid)
			}
		})
	}
}

func TestDecodeGameGrid(t *testing.T) {
	testCases := []struct {
		name          string
		encodedString string
		gridSize      int
		expected      [][]Cell
	}{
		{
			name:          "Simple decode test",
			encodedString: "MR|EM|",
			gridSize:      2,
			expected: [][]Cell{
				{{HasMine: true, AdjacentMines: 1}, {IsRevealed: true, AdjacentMines: 2}},
				{{AdjacentMines: 2}, {HasMine: true, AdjacentMines: 1}},
			},
		},
		{
			name:          "All cells flagged with mines",
			encodedString: "XX|XX|",
			gridSize:      2,
			expected: [][]Cell{
				{{IsFlagged: true, HasMine: true, AdjacentMines: 3}, {IsFlagged: true, HasMine: true, AdjacentMines: 3}},
				{{IsFlagged: true, HasMine: true, AdjacentMines: 3}, {IsFlagged: true, HasMine: true, AdjacentMines: 3}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decodedGameGrid := DecodeGameGrid(tc.encodedString, tc.gridSize)
			if !reflect.DeepEqual(decodedGameGrid, tc.expected) {
				t.Errorf("Test case '%s' failed. Expected decoded grid to be '%v', but got '%v'", tc.name, tc.expected, decodedGameGrid)
			}
		})
	}
}
