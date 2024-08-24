package internal

import (
	"fmt"
	"html/template"
	"minesweeper/src"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

type Handler struct {
	Templates *template.Template
	Store     *sessions.CookieStore
}

func NewHandler(templates *template.Template, store *sessions.CookieStore) *Handler {
	return &Handler{Templates: templates, Store: store}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {

	err := h.Templates.ExecuteTemplate(w, "index", nil)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) StartGame(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to process form", http.StatusBadRequest)
		return
	}

	// get the form data
	gridSize, gridSizeErr := strconv.Atoi(r.FormValue("grid-size"))
	minesAmount, minesAmountErr := strconv.Atoi((r.FormValue("mines-amount")))

	if gridSizeErr != nil || minesAmountErr != nil {
		http.Error(w, "Invalid input values", http.StatusBadRequest)
		return
	}

	if gridSize <= 0 || minesAmount <= 0 {
		http.Error(w, "The grid size and mines amount must be greater than 0", http.StatusUnprocessableEntity)
		return
	}

	game := src.NewGame(gridSize, minesAmount)

	if err := SaveGameToSession(w, r, game, h.Store); err != nil {
		http.Error(w, fmt.Sprintf("Error saving game to session: %v", err), http.StatusInternalServerError)
		return
	}

	gameGridHtml, gridGenerationErr := GenerateGridHTML(h.Templates, game)
	if gridGenerationErr != nil {
		http.Error(w, fmt.Sprintf("Error generating grid HTML: %v", gridGenerationErr), http.StatusInternalServerError)
		return
	}

	responseData := struct {
		GridSize     int
		MinesAmount  int
		GameGridHtml template.HTML
	}{
		GridSize:     gridSize,
		MinesAmount:  minesAmount,
		GameGridHtml: template.HTML(gameGridHtml),
	}

	err := h.Templates.ExecuteTemplate(w, "game_layout", responseData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RevealCell(w http.ResponseWriter, r *http.Request) {
	game, err := GetGameFromSession(r, h.Store)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get game from session: %v", err), http.StatusInternalServerError)
		return
	}

	rowStr := r.URL.Query().Get("row")
	colStr := r.URL.Query().Get("col")

	row, err := strconv.Atoi(rowStr)
	if err != nil {
		http.Error(w, "Invalid row value", http.StatusBadRequest)
		return
	}

	col, err := strconv.Atoi(colStr)
	if err != nil {
		http.Error(w, "Invalid column value", http.StatusBadRequest)
		return
	}

	game.RevealCell(row, col)

	if err := SaveGameToSession(w, r, game, h.Store); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save game: %v", err), http.StatusInternalServerError)
		return
	}

	gameGridHtml, gridGenerationErr := GenerateGridHTML(h.Templates, game)
	if gridGenerationErr != nil {
		http.Error(w, fmt.Sprintf("Error generating grid HTML: %v", gridGenerationErr), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(gameGridHtml))
}

func (h *Handler) FlagCell(w http.ResponseWriter, r *http.Request) {
	row, _ := strconv.Atoi(r.URL.Query().Get("row"))
	col, _ := strconv.Atoi(r.URL.Query().Get("col"))

	game, _ := GetGameFromSession(r, h.Store)

	game.FlagCell(row, col)

	if game.CheckWinCondition() {
		game.GameWon = true
	}

	if err := SaveGameToSession(w, r, game, h.Store); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save game: %v", err), http.StatusInternalServerError)
		return
	}

	gameGridHtml, _ := GenerateGridHTML(h.Templates, game)

	w.Write([]byte(gameGridHtml))
}
