package main

import (
	"cmp"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"minesweeper/internal"
	"minesweeper/internal/db"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

// Package-level variable to store the parsed templates
var templates *template.Template

var globalStore *sessions.CookieStore

func init() {
	var err error

	// Parse templates from the "public" and "public/**" directories
	funcMap := template.FuncMap{
		"Sub": func(a int, b int) int { return a - b },
		"Add": func(a int, b int) int { return a + b },
	}

	templates, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	templates, err = templates.New("").Funcs(funcMap).ParseGlob("templates/**/*.html")
	if err != nil {
		log.Fatalf("Error parsing components: %v", err)
	}

	// Log the templates that have been parsed
	for _, tmpl := range templates.Templates() {
		log.Printf("Parsed template: %s", tmpl.Name())
	}

	fmt.Println("Trying to load .env file...")
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Println(".env file not found. Continuing with environment variables.")
	}

	// Now load the SESSION_SECRET from environment variables
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET environment variable not set. The application cannot start without it.")
	}

	globalStore = sessions.NewCookieStore([]byte(sessionSecret))
	// TODO fix this so it will also work with HTTPS
	globalStore.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,  // seconds
		Secure:   false, // Set to true if using HTTPS
	}
	log.Println("Session store initialized successfully.")
}

func connectToDB() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	db, err := sql.Open("sqlite", databaseURL)
	if err != nil {
		return nil, err
	}

	// Check if the database is accessible
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	logFile, err := os.OpenFile("minesweeper.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	mux := http.NewServeMux()

	staticFiles := http.StripPrefix("/dist/", http.FileServer(http.Dir("dist")))
	mux.Handle("/dist/", staticFiles)

	dbConn, err := connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)
	handler := internal.NewHandler(templates, globalStore, queries)

	mux.HandleFunc("/", handler.Index)
	mux.HandleFunc("/load-game", handler.LoadGame)
	mux.HandleFunc("/start-game", handler.StartGame)
	mux.HandleFunc("/handle-grid-action", handler.HandleGridAction)
	mux.HandleFunc("/games", handler.IndexGames)
	mux.HandleFunc("/session-games-info", handler.SessionGamesInfo)

	port := cmp.Or(os.Getenv("APP_PORT"), "8080")

	fmt.Printf("Server is listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
