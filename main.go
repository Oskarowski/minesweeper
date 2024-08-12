package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func main() {
	http.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("dist"))))

	templates, err := template.ParseFiles("public/index.html", "public/name_card.html")

	if err != nil {
		panic(fmt.Sprintf("Error parsing template: %v", err))
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res := map[string]interface{}{
			"Name":  "Wyndham",
			"Phone": "8888888",
			"Email": "skyscraper@gmail.com",
		}

		err := templates.ExecuteTemplate(w, "index", res)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
