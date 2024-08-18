package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"minesweeper/src"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

// Package-level variable to store the parsed templates
var templates *template.Template

var globalStore *sessions.CookieStore

func init() {
	var err error

	// Parse templates from the "public" and "public/components" directories
	templates, err = template.ParseGlob("public/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	templates, err = templates.ParseGlob("public/components/*.html")
	if err != nil {
		log.Fatalf("Error parsing components: %v", err)
	}

	// Log the templates that have been parsed
	for _, tmpl := range templates.Templates() {
		log.Printf("Parsed template: %s", tmpl.Name())
	}

	gob.Register(&src.Game{})
	gob.Register(&src.Cell{})

	fmt.Println("Trying to load .env file...")
	envErr := godotenv.Load(".env")
	if envErr != nil {
		panic(envErr)
	}

	globalStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
}

func saveGameToSession(w http.ResponseWriter, r *http.Request, game *src.Game) error {
	session, err := globalStore.Get(r, "minesweeper-session")

	if err != nil {
		return err
	}

	session.Values["game"] = game
	return session.Save(r, w)
}

func getGameFromSession(r *http.Request) (*src.Game, error) {
	session, err := globalStore.Get(r, "minesweeper-session")

	if err != nil {
		return nil, err
	}

	if game, ok := session.Values["game"].(*src.Game); ok {
		return game, nil
	}

	return nil, fmt.Errorf("game not found in session")
}

func startGameHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := saveGameToSession(w, r, game); err != nil {
		http.Error(w, fmt.Sprintf("Error saving game to session: %v", err), http.StatusInternalServerError)
		return
	}

	gameGridHtml := src.GenerateGridHTML(game)

	responseData := struct {
		GridSize     int
		MinesAmount  int
		GameGridHtml template.HTML
	}{
		GridSize:     gridSize,
		MinesAmount:  minesAmount,
		GameGridHtml: template.HTML(gameGridHtml),
	}

	err := templates.ExecuteTemplate(w, "game_grid", responseData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

func revealCellHandler(w http.ResponseWriter, r *http.Request) {
	game, err := getGameFromSession(r)

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

	if err := saveGameToSession(w, r, game); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save game: %v", err), http.StatusInternalServerError)
		return
	}

	gameGridHtml := src.GenerateGridHTML(game)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(gameGridHtml))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index", nil)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

	http.HandleFunc("/", indexHandler)

	http.HandleFunc("/start-game", startGameHandler)

	http.HandleFunc("/reveal", revealCellHandler)

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
