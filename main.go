package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type UserInfo struct {
	Name  string
	Phone string
	Email string
}

type Button struct {
	Action      string
	Target      string
	Label       string
	ButtonClass string
	HoverClass  string
}

func main() {
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

	templates, err := template.ParseGlob("public/*.html")

	if err != nil {
		panic(fmt.Sprintf("Error parsing template: %v", err))
	}

	templates, err = templates.ParseGlob("public/components/*.html")
	if err != nil {
		panic(fmt.Sprintf("Error parsing components: %v", err))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		revealUserInfoBtn := Button{
			Action:      "/get-user-info",
			Target:      "#user-info-container",
			Label:       "Reveal the User Info",
			ButtonClass: "bg-blue-600",
			HoverClass:  "bg-blue-700",
		}

		data := struct {
			RevealUserButton Button
		}{
			RevealUserButton: revealUserInfoBtn,
		}

		err := templates.ExecuteTemplate(w, "index", data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/get-user-info", func(w http.ResponseWriter, r *http.Request) {
		userInfo := UserInfo{
			Name:  "Wyndham",
			Phone: "8888888",
			Email: "skyscraper@gmail.com",
		}

		button := Button{
			Action:      "/hide-user-info",
			Target:      "#user-info-container",
			Label:       "Hide User Info",
			ButtonClass: "bg-red-600",
			HoverClass:  "bg-red-700",
		}

		data := struct {
			UserInfo       UserInfo
			HideButtonData Button
		}{
			UserInfo:       userInfo,
			HideButtonData: button,
		}

		err := templates.ExecuteTemplate(w, "user_info_card", data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/hide-user-info", func(w http.ResponseWriter, r *http.Request) {
		button := Button{
			Action:      "/get-user-info",
			Target:      "#user-info-container",
			Label:       "Reveal the User Info",
			ButtonClass: "bg-blue-600",
			HoverClass:  "bg-blue-700",
		}

		err := templates.ExecuteTemplate(w, "button", button)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering button: %v", err), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
