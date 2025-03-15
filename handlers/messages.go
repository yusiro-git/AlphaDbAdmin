package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
)

type MessageHandler struct {
	db *pgx.Conn
}

func NewMessageHandler(conn *pgx.Conn) MessageHandler {
	return MessageHandler{db: conn}
}

func (handler MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	rows, err := handler.db.Query(context.Background(), "SELECT id, message, date FROM delayedmessages")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var delayedmesages []map[string]string
	for rows.Next() {
		var id int
		var message string
		var date time.Time
		err = rows.Scan(&id, &message, &date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		delayedmesages = append(delayedmesages, map[string]string{
			"id":      strconv.Itoa(id),
			"message": message,
			"date":    date.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delayedmesages)
}

func (handler MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := handler.db.Exec(context.Background(), "DELETE FROM delayedmassages WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	date := r.FormValue("date")
	_, err := handler.db.Exec(context.Background(), "INSERT INTO delayedmassage (message, date) VALUES ($1, $2)", message, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (handler MessageHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	rows, err := handler.db.Query(context.Background(), "SELECT id, message, date FROM delayedmessages WHERE message ILIKE $1", "%"+query+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var usefullinkss []map[string]string
	for rows.Next() {
		var id int
		var message string
		var date time.Time
		err = rows.Scan(&id, &message, &date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		usefullinkss = append(usefullinkss, map[string]string{
			"id":      strconv.Itoa(id),
			"message": message,
			"date":    date.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usefullinkss)
}

func (handler MessageHandler) HtmlMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering messages page")
	tmpl := template.Must(template.ParseFiles("templates/messages.html"))
	tmpl.Execute(w, nil)
}
