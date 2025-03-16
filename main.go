package main

import (
	"AlphaDbAdmin/handlers"
	"AlphaDbAdmin/middleware"
	"log"
	"net/http"

	"AlphaDbAdmin/storage/postgres"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	db, err := postgres.New("postgres://postgres:Go@localhost:5432/alphatgbot")
	if err != nil {
		panic(err)
	}

	// Router setup
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/auth", handlers.AuthHandler).Methods("GET", "POST")

	faq := handlers.NewFAQHandler(db)
	links := handlers.NewLinksHandler(db)
	messages := handlers.NewMessageHandler(db)

	// Protected routes
	log.Println("Protected routes")
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/", handlers.HomeHandler).Methods("GET") // Homepage route
	protected.HandleFunc("/api/faqs", faq.HtmlFAQ)
	protected.HandleFunc("/api/faqs/load", faq.GetFAQs).Methods("GET")
	protected.HandleFunc("/api/faqs/delete", faq.DeleteFAQ).Methods("DELETE")
	protected.HandleFunc("/api/faqs/create", faq.CreateFAQ).Methods("POST")
	protected.HandleFunc("/api/faqs/search", faq.SearchFAQ).Methods("GET")

	protected.HandleFunc("/api/links", links.HtmlLinks)
	protected.HandleFunc("/api/links/load", links.GetLinks).Methods("GET")
	protected.HandleFunc("/api/links/delete", links.DeleteLink).Methods("DELETE")
	protected.HandleFunc("/api/links/create", links.CreateLink).Methods("POST")
	protected.HandleFunc("/api/links/search", links.SearchLinks).Methods("GET")

	protected.HandleFunc("/api/messages", messages.HtmlMessages)
	protected.HandleFunc("/api/messages/load", messages.GetMessages).Methods("GET")
	protected.HandleFunc("/api/messages/delete", messages.DeleteMessage).Methods("DELETE")
	protected.HandleFunc("/api/messages/create", messages.CreateMessage).Methods("POST")
	protected.HandleFunc("/api/messages/search", messages.SearchMessages).Methods("GET")
	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start server
	log.Println("Server started on :8081")
	log.Fatal(http.ListenAndServe("0.0.0.0:8081", r))
}
