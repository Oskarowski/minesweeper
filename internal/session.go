package internal

import (
	"fmt"
	"minesweeper/src"
	"net/http"

	"github.com/gorilla/sessions"
)

func SaveGameToSession(w http.ResponseWriter, r *http.Request, game *src.Game, store *sessions.CookieStore) error {
	session, err := store.Get(r, "minesweeper-session")

	if err != nil {
		return err
	}

	session.Values["game"] = game
	return session.Save(r, w)
}

func GetGameFromSession(r *http.Request, store *sessions.CookieStore) (*src.Game, error) {
	session, err := store.Get(r, "minesweeper-session")

	if err != nil {
		return nil, err
	}

	if game, ok := session.Values["game"].(*src.Game); ok {
		return game, nil
	}

	return nil, fmt.Errorf("game not found in session")
}
