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

type FAQHandler struct {
	db *pgx.Conn
}

func NewFAQHandler(conn *pgx.Conn) FAQHandler {
	return FAQHandler{db: conn}
}

func (handler FAQHandler) GetFAQs(w http.ResponseWriter, r *http.Request) {
	rows, err := handler.db.Query(context.Background(), "SELECT id, question, answer FROM FAQ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var faqs []map[string]string
	for rows.Next() {
		var id int
		var question, answer string
		err = rows.Scan(&id, &question, &answer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		faqs = append(faqs, map[string]string{
			"id":       strconv.Itoa(id),
			"question": question,
			"answer":   answer,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faqs)
}

func (handler FAQHandler) DeleteFAQ(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := handler.db.Exec(context.Background(), "DELETE FROM FAQ WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler FAQHandler) CreateFAQ(w http.ResponseWriter, r *http.Request) {
	question := r.FormValue("question")
	answer := r.FormValue("answer")
	_, err := handler.db.Exec(context.Background(), "INSERT INTO FAQ (question, answer) VALUES ($1, $2)", question, answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (handler FAQHandler) SearchFAQ(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	rows, err := handler.db.Query(context.Background(), "SELECT id, question, answer FROM FAQ WHERE question ILIKE $1", "%"+query+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var faqs []map[string]string
	for rows.Next() {
		var id int
		var question, answer string
		err = rows.Scan(&id, &question, &answer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		faqs = append(faqs, map[string]string{
			"id":       strconv.Itoa(id),
			"question": question,
			"answer":   answer,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(faqs)
}

func (handler FAQHandler) HtmlFAQ(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering messages page")
	tmpl := template.Must(template.ParseFiles("templates/faq.html"))
	tmpl.Execute(w, nil)
}
