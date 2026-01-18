package main 

import (
	"database/sql"
	"encoding/json"
    "net/http"
	"time"
	"log"
	_ "github.com/mattn/go-sqlite3"

)

type Topic struct {
	ID          int       `json:"id"`
	Title        string   `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
}

var db *sql.DB

func main() {
	var err error

	db, err = sql.Open("sqlite3", "../.db/magic.db")
	if err != nil {
		log.Fatal("Failed to open database", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database", err)
	}

	log.Println("Connected to database!")

	http.HandleFunc("/topics", topicHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server", err)
	}
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        handleGetTopic(w, r)
    case http.MethodPost:
        handlePostTopic(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func handleGetTopic(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT ID, Title, Description, Created_At, Updated_At FROM topics
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var topics []Topic

	for rows.Next() {
		var topic Topic
		if err := rows.Scan(
			&topic.ID,
			&topic.Title,
			&topic.Description,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		topics = append(topics, topic)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

func handlePostTopic(w http.ResponseWriter, r *http.Request) {
    var topic Topic
    err := json.NewDecoder(r.Body).Decode(&topic)
    if err != nil {
        http.Error(w, `{"error":"invalid JSON: `+err.Error()+`"}`, http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(topic)
}