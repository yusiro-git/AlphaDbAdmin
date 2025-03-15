package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v4"
)

type LinksHandler struct {
	db *pgx.Conn
}

func NewLinksHandler(conn *pgx.Conn) LinksHandler {
	return LinksHandler{db: conn}
}

func (handler LinksHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	rows, err := handler.db.Query(context.Background(), "SELECT id, linktext, url FROM usefullinks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var usefullinkss []map[string]string
	for rows.Next() {
		var id int
		var linktext, url string
		err = rows.Scan(&id, &linktext, &url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		usefullinkss = append(usefullinkss, map[string]string{
			"id":       strconv.Itoa(id),
			"linktext": linktext,
			"url":      url,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usefullinkss)
}

func (handler LinksHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := handler.db.Exec(context.Background(), "DELETE FROM usefullinks WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler LinksHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	question := r.FormValue("linktext")
	answer := r.FormValue("url")
	_, err := handler.db.Exec(context.Background(), "INSERT INTO usefullinks (linktext, url) VALUES ($1, $2)", question, answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (handler LinksHandler) SearchLinks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	rows, err := handler.db.Query(context.Background(), "SELECT id, linktext, url FROM usefullinks WHERE linktext ILIKE $1", "%"+query+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var usefullinkss []map[string]string
	for rows.Next() {
		var id int
		var linktext, url string
		err = rows.Scan(&id, &linktext, &url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		usefullinkss = append(usefullinkss, map[string]string{
			"id":       strconv.Itoa(id),
			"linktext": linktext,
			"url":      url,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usefullinkss)
}

func (handler LinksHandler) HtmlLinks(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering messages page")
	tmpl := template.Must(template.ParseFiles("templates/links.html"))
	tmpl.Execute(w, nil)
}
