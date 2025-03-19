package main

import (
	"AlphaDbAdmin/config"
	"AlphaDbAdmin/handlers"
	"AlphaDbAdmin/middleware"
	"AlphaDbAdmin/storage/postgres"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	conf := config.MustLoad()

	//postgres://postgres:Go@localhost:5432/alphatgbot
	db, err := postgres.New(conf.PostresConnectionURL)
	if err != nil {
		panic(err)
	}

	// Router setup
	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/auth", handlers.AuthHandler(conf.AdminKey)).Methods("GET", "POST")

	faq := handlers.NewFAQHandler(db)
	links := handlers.NewLinksHandler(db)
	messages := handlers.NewMessageHandler(db)

	// Protected routes
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
	log.Println("Server started on :6969")
	log.Fatal(http.ListenAndServe("0.0.0.0:6969", r))
}
