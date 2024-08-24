package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"minesweeper/internal"
	"minesweeper/src"
	"net/http"
	"os"

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

	gob.Register(&src.Game{})
	gob.Register(&src.Cell{})

	fmt.Println("Trying to load .env file...")
	envErr := godotenv.Load(".env")
	if envErr != nil {
		panic(envErr)
	}

	globalStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
}

func main() {
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

	handler := internal.NewHandler(templates, globalStore)

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
