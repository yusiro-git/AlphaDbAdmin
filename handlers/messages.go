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
	rows, err := handler.db.Query(context.Background(), "SELECT id, group_name, message, date, image FROM delayedmessages")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var delayedMessages []map[string]string
	for rows.Next() {
		var id int
		var groupName, message, pictureUrl string
		var date time.Time
		err = rows.Scan(&id, &groupName, &message, &date, &pictureUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		delayedMessages = append(delayedMessages, map[string]string{
			"id":          strconv.Itoa(id),
			"group_name":  groupName,
			"message":     message,
			"date":        date.Format("2006-01-02 15:04:05"),
			"picture_url": pictureUrl,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delayedMessages)
}

func (handler MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := handler.db.Exec(context.Background(), "DELETE FROM delayedmessages WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	groupName := r.FormValue("group_name")
	message := r.FormValue("message")
	date := r.FormValue("date")
	pictureUrl := r.FormValue("picture_url")

	_, err := handler.db.Exec(context.Background(), "INSERT INTO delayedmessages (group_name, message, date, image) VALUES ($1, $2, $3, $4)", groupName, message, date, pictureUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (handler MessageHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	rows, err := handler.db.Query(context.Background(), "SELECT id, group_name, message, date, image FROM delayedmessages WHERE message ILIKE $1 OR group_name ILIKE $1", "%"+query+"%")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var usefulLinks []map[string]string
	for rows.Next() {
		var id int
		var groupName, message, pictureUrl string
		var date time.Time
		err = rows.Scan(&id, &groupName, &message, &date, &pictureUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		usefulLinks = append(usefulLinks, map[string]string{
			"id":          strconv.Itoa(id),
			"group_name":  groupName,
			"message":     message,
			"date":        date.Format("2006-01-02 15:04:05"),
			"picture_url": pictureUrl,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usefulLinks)
}

func (handler MessageHandler) HtmlMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering messages page")
	tmpl := template.Must(template.ParseFiles("templates/messages.html"))
	tmpl.Execute(w, nil)
}
