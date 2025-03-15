package handlers

import (
	"net/http"
	"text/template"
)

// HomeHandler handles requests to the homepage ("/")
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authorized (middleware already handles this)
	// Render the homepage template
	if r.RequestURI != "/" {
		return
	}
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Failed to load homepage template", http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render homepage", http.StatusInternalServerError)
		return
	}
}
