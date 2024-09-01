package internal

import (
	"fmt"
	"log"
	"minesweeper/internal/models"
	"net/http"

	"github.com/gorilla/sessions"
)

func SaveGameToSession(w http.ResponseWriter, r *http.Request, game *models.Game, store *sessions.CookieStore) error {
	if game.Uuid == "" {
		return fmt.Errorf("game uuid is empty")
	}

	session, err := store.Get(r, "minesweeper-session")

	if err != nil {
		log.Printf("Failed to get session: %v", err)
		return fmt.Errorf("failed to get session: %w", err)
	}

	var uuids []string

	// Retrieve existing UUIDs from the session, if available
	if storedUuids, ok := session.Values["game_uuids"]; ok {
		// Check type assertion
		if uuids, ok = storedUuids.([]string); !ok {
			log.Printf("Unexpected type for game_uuids: %T. Resetting UUIDs.", storedUuids)
			uuids = []string{}
		} else {
			log.Printf("Retrieved existing game UUIDs from session: %v", uuids)
		}
	} else {
		log.Printf("No existing game UUIDs found in session.")
	}

	if !contains(uuids, game.Uuid) {
		uuids = append(uuids, game.Uuid)
		session.Values["game_uuids"] = uuids
		log.Printf("Added new game UUID to session: %s", game.Uuid)
	} else {
		log.Printf("Game UUID %s already exists in session, not adding.", game.Uuid)
	}

	if err := session.Save(r, w); err != nil {
		log.Printf("Failed to save session: %v", err)
		return fmt.Errorf("failed to save session: %w", err)
	}

	log.Printf("Session saved successfully. Elements in slice: %d Updated game UUIDs in session: %+q", len(uuids), uuids)
	return nil
}

func GetGameFromSession(r *http.Request, store *sessions.CookieStore) ([]string, error) {
	session, err := store.Get(r, "minesweeper-session")
	if err != nil {
		return nil, err
	}

	rawUuids, exists := session.Values["game_uuids"]
	if !exists {
		return nil, fmt.Errorf("no game UUIDs found in session")
	}

	uuids, ok := rawUuids.([]string)
	if !ok {
		return nil, fmt.Errorf("game UUIDs in session are not of type []string, but got %T", rawUuids)
	}

	return uuids, nil
}

func contains(s []string, str string) bool {
	if len(s) == 0 {
		return false
	}

	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
