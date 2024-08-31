package main

import (
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
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	templates, err = templates.ParseGlob("templates/**/*.html")
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
		panic(envErr)
	}

	globalStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
}

func main() {
	logFile, err := os.OpenFile("minesweeper.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

	databaseURL := os.Getenv("DATABASE_URL")
	dbConn, err := sql.Open("sqlite", databaseURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)

	handler := internal.NewHandler(templates, globalStore, queries)

	http.HandleFunc("/", handler.Index)

	http.HandleFunc("/start-game", handler.StartGame)

	http.HandleFunc("/reveal", handler.RevealCell)

	http.HandleFunc("/flag", handler.FlagCell)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
