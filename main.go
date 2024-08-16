package main

import (
	"fmt"
	"html/template"
	"log"
	"minesweeper/src"
	"net/http"
	"strconv"
)

// Package-level variable to store the parsed templates
var templates *template.Template

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

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
