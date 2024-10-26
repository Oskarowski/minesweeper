package internal

import (
	"fmt"
	"html/template"
	"log"
	"minesweeper/internal/db"
	"minesweeper/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

type Handler struct {
	Templates *template.Template
	Store     *sessions.CookieStore
	Queries   *db.Queries
}

func NewHandler(templates *template.Template, store *sessions.CookieStore, queries *db.Queries) *Handler {
	return &Handler{Templates: templates, Store: store, Queries: queries}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	gameUuid := r.URL.Query().Get("game_uuid")

	responseData := map[string]interface{}{
		"HasGameUuid": gameUuid != "",
		"GameUuid":    gameUuid,
	}

	err := h.Templates.ExecuteTemplate(w, "index", responseData)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LoadGame(w http.ResponseWriter, r *http.Request) {

	gameUuid := r.URL.Query().Get("game_uuid")
	if gameUuid == "" {
		http.Error(w, "Missing game-uuid parameter", http.StatusBadRequest)
		return
	}

	dbGame, dbGameRetrieveErr := h.Queries.GetGameByUuid(r.Context(), gameUuid)
	if dbGameRetrieveErr != nil {
		http.Error(w, fmt.Sprintf("Not able to retrieve game with such uuid : %v", dbGameRetrieveErr), http.StatusConflict)
		return
	}

	game, gameModelErr := models.FromDbGame(&dbGame)
	if gameModelErr != nil {
		http.Error(w, fmt.Sprintf("Error during game casting: %v", gameModelErr), http.StatusInternalServerError)
		return
	}

	if err := SaveGameToSession(w, r, game, h.Store); err != nil {
		http.Error(w, fmt.Sprintf("Not able to save game to session: %v", err), http.StatusInternalServerError)
		return
	}

	gameGridHtml, gridGenerationErr := GenerateGridHTML(h.Templates, game)
	if gridGenerationErr != nil {
		http.Error(w, fmt.Sprintf("Not able to generate game grid: %v", gridGenerationErr), http.StatusInternalServerError)
		return
	}

	responseData := struct {
		GridSize     int
		MinesAmount  int
		GameGridHtml template.HTML
	}{
		GridSize:     int(dbGame.GridSize),
		MinesAmount:  int(dbGame.MinesAmount),
		GameGridHtml: template.HTML(gameGridHtml),
	}

	err := h.Templates.ExecuteTemplate(w, "game_layout", responseData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) StartGame(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.returnErrorResponse(ErrorResponseConfig{
			ResponseWriter: w,
			ErrorMessage:   err.Error(),
			ShowCloseBtn:   true,
		})
		log.Printf("Error parsing form: %v", err)
		return
	}

	gameSettings, formValidationErr := ValidateGameSettingsForm(
		r.FormValue("grid-size"),
		r.FormValue("mines-amount"),
		r.FormValue("random-mines"),
		r.FormValue("random-grid-size"),
	)

	if formValidationErr != nil {
		h.returnErrorResponse(ErrorResponseConfig{
			ResponseWriter: w,
			ErrorMessage:   formValidationErr.Error(),
			ShowCloseBtn:   true,
		})
		log.Printf("Form validation error: %v", formValidationErr)
		return
	}

	// TODO eliminate need for this double creation of game because how grid state is initialized
	newGame := models.NewGame(gameSettings.GridSize, gameSettings.MinesAmount)
	dbGame, dbGameErr := h.Queries.CreateGame(r.Context(), db.CreateGameParams{
		GridSize:    int64(gameSettings.GridSize),
		MinesAmount: int64(gameSettings.MinesAmount),
		GridState:   models.EncodeGameGrid(newGame.Grid),
	})

	if dbGameErr != nil {
		h.returnErrorResponse(ErrorResponseConfig{
			ResponseWriter: w,
			ErrorMessage:   fmt.Sprintf("Error creating game: %v", dbGameErr),
			ShowCloseBtn:   false,
		})
		return
	}

	game, gameModelErr := models.FromDbGame(&dbGame)

	if gameModelErr != nil {
		h.returnErrorResponse(ErrorResponseConfig{
			ResponseWriter: w,
			ErrorMessage:   fmt.Sprintf("Error creating game model: %v", gameModelErr),
			ShowCloseBtn:   false,
		})
		return
	}

	if err := SaveGameToSession(w, r, game, h.Store); err != nil {
		h.returnErrorResponse(ErrorResponseConfig{
			ResponseWriter: w,
			ErrorMessage:   fmt.Sprintf("Error saving game to session: %v", err),
			ShowCloseBtn:   false,
		})
		return
	}

	gameGridHtml, gridGenerationErr := GenerateGridHTML(h.Templates, game)
	if gridGenerationErr != nil {
		h.returnErrorResponse(ErrorResponseConfig{
			ResponseWriter: w,
			ErrorMessage:   fmt.Sprintf("Error generating grid HTML: %v", gridGenerationErr),
			ShowCloseBtn:   false,
		})
		return
	}

	responseData := struct {
		GridSize     int
		MinesAmount  int
		GameGridHtml template.HTML
	}{
		GridSize:     int(dbGame.GridSize),
		MinesAmount:  int(dbGame.MinesAmount),
		GameGridHtml: template.HTML(gameGridHtml),
	}

	w.Header().Set("HX-Trigger", "gameStarted")
	err := h.Templates.ExecuteTemplate(w, "game_layout", responseData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleGridAction(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	row, rowParErr := strconv.Atoi(r.URL.Query().Get("row"))
	col, colParErr := strconv.Atoi(r.URL.Query().Get("col"))

	if action == "" || rowParErr != nil || colParErr != nil {
		http.Error(w, "Unprocessable or missing request parameters.", http.StatusUnprocessableEntity)
		return
	}

	storedGamesUuids, err := GetGameFromSession(r, h.Store)

	if err != nil || len(storedGamesUuids) == 0 {
		log.Printf("Failed to get game from session: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get game from session: %v", err), http.StatusInternalServerError)
		return
	}
	lastGameUuid := storedGamesUuids[len(storedGamesUuids)-1]

	dbGame, err := h.Queries.GetGameByUuid(r.Context(), lastGameUuid)
	if err != nil {
		log.Printf("Failed to get game from database: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get game from database: %v", err), http.StatusInternalServerError)
		return
	}

	game, err := models.FromDbGame(&dbGame)
	if err != nil {
		log.Printf("Failed to convert db game: %v", err)
		http.Error(w, fmt.Sprintf("Failed to convert db game: %v", err), http.StatusInternalServerError)
		return
	}

	switch action {
	case "reveal_cell":
		game.RevealCell(row, col)
	case "flag_cell":
		game.FlagCell(row, col)
	default:
		http.Error(w, "Invalid action.", http.StatusBadRequest)
		return
	}

	encodedGridState := models.EncodeGameGrid(game.Grid)
	err = h.Queries.UpdateGameGridStateById(r.Context(), db.UpdateGameGridStateByIdParams{
		GameFailed: game.GameFailed,
		GameWon:    game.GameWon,
		GridState:  encodedGridState,
		Id:         game.Id,
	})
	if err != nil {
		log.Printf("Failed to update game state in database: %v", err)
		http.Error(w, fmt.Sprintf("Failed to update game state in database: %v", err), http.StatusInternalServerError)
		return
	}

	gameGridHtml, gridGenerationErr := GenerateGridHTML(h.Templates, game)
	if gridGenerationErr != nil {
		http.Error(w, fmt.Sprintf("Error generating grid HTML: %v", gridGenerationErr), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(gameGridHtml))
}

func (h *Handler) IndexGames(w http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("page")
	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	var pageSize int64 = 25

	offset := (pageNumber - 1) * int(pageSize) // pageSize

	games, err := h.Queries.ListGames(r.Context(), db.ListGamesParams{
		Limit:  pageSize,
		Offset: int64(offset),
	})

	if err != nil {
		log.Printf("Failed to get games from database: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get games from database: %v", err), http.StatusInternalServerError)
		return
	}

	totalGamesCount, totalGamesCountErr := GetTotalGamesCount(h.Queries)
	if totalGamesCountErr != nil {
		log.Printf("Failed to get total games count: %v", totalGamesCountErr)
		http.Error(w, fmt.Sprintf("Failed to get total games count: %v", totalGamesCountErr), http.StatusInternalServerError)
		return
	}

	totalPages := totalGamesCount / pageSize

	data := struct {
		Games           []db.ListGamesRow
		CurrentPage     int
		TotalPages      int
		TotalGamesCount int
	}{
		Games:           games,
		CurrentPage:     pageNumber,
		TotalPages:      int(totalPages),
		TotalGamesCount: int(totalGamesCount),
	}

	if err := h.Templates.ExecuteTemplate(w, "index_games_page", data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, fmt.Sprintf("Failed to render template: %v", err), http.StatusInternalServerError)
		return
	}
}

type ErrorResponseConfig struct {
	ResponseWriter http.ResponseWriter
	ErrorMessage   string
	ShowCloseBtn   bool
}

func (h *Handler) returnErrorResponse(config ErrorResponseConfig) {
	if config.ErrorMessage == "" {
		config.ErrorMessage = "An error occurred"
	}

	responseData := struct {
		ErrorMessage string
		ShowCloseBtn bool
	}{
		ErrorMessage: config.ErrorMessage,
		ShowCloseBtn: config.ShowCloseBtn,
	}

	config.ResponseWriter.WriteHeader(http.StatusBadRequest)
	err := h.Templates.ExecuteTemplate(config.ResponseWriter, "error_message", responseData)
	if err != nil {
		http.Error(config.ResponseWriter, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SessionGamesInfo(w http.ResponseWriter, r *http.Request) {
	storedUuids, sessionErr := GetGameFromSession(r, h.Store)

	if len(storedUuids) == 0 || sessionErr != nil {
		err := h.Templates.ExecuteTemplate(w, "session_games_info", nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		}
		return
	}

	gamesInSessionInfo, err := h.Queries.GetGamesInfoByUuids(r.Context(), storedUuids)

	if err != nil {
		log.Printf("Failed to fetch game counts from database: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch game counts: %v", err), http.StatusInternalServerError)
		return
	}

	responseData := struct {
		HasGames         bool
		TotalGames       int
		WonGames         int
		LostGames        int
		NotFinishedGames int
	}{
		HasGames:         gamesInSessionInfo.TotalGames > 0,
		TotalGames:       int(gamesInSessionInfo.TotalGames),
		LostGames:        int(gamesInSessionInfo.LostGames),
		WonGames:         int(gamesInSessionInfo.WonGames),
		NotFinishedGames: int(gamesInSessionInfo.NotFinishedGames),
	}

	err = h.Templates.ExecuteTemplate(w, "session_games_info", responseData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) Charts(w http.ResponseWriter, r *http.Request) {
	err := h.Templates.ExecuteTemplate(w, "charts_page", nil)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}
