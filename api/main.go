package main 

import (
	"database/sql"
	"encoding/json"
    "net/http"
	"time"
	"log"
	"strings"
	"strconv"
	_ "github.com/mattn/go-sqlite3"

)
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type Topic struct {
	ID          int       `json:"id"`
	Title        string   `json:"title"`
	Description string    `json:"description"`
	UserID    int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    string    `json:"user_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
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

	http.HandleFunc("/users", userHandler)
	http.HandleFunc("/topics", topicHandler)
	http.HandleFunc("/posts", postHandler)
	http.HandleFunc("/comments", commentHandler)

	http.HandleFunc("/users/", handleGetUserByID)
	http.HandleFunc("/topics/", handleGetTopicByID)
	http.HandleFunc("/posts/", handleGetPostByID)
	http.HandleFunc("/comments/", handleGetCommentByID)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server", err)
	}
}

//USERS
func userHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        handleGetUser(w, r)
    case http.MethodPost:
        handlePostUser(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT ID, Username, Created_At FROM users
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.CreatedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func handlePostUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"error":"invalid JSON: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/users/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    row := db.QueryRow(`
        SELECT id, username, created_at
        FROM users
        WHERE id = ?
    `, id)

    var user User
    err = row.Scan(
        &user.ID,
        &user.Username,
        &user.CreatedAt,
    )

    if err == sql.ErrNoRows {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

//TOPICS
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
        SELECT ID, Title, Description, User_ID, Created_At, Updated_At FROM topics
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
			&topic.UserID,
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

func handleGetTopicByID(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/topics/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    row := db.QueryRow(`
        SELECT id, title, description, user_id, created_at, updated_at
        FROM topics
        WHERE id = ?
    `, id)

    var topic Topic
    err = row.Scan(
        &topic.ID,
        &topic.Title,
        &topic.Description,
        &topic.UserID,
        &topic.CreatedAt,
        &topic.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        http.Error(w, "Topic not found", http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(topic)
}

//POSTS
func postHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        handleGetPost(w, r)
    case http.MethodPost:
        handlePostPost(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func handleGetPost(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT ID, User_ID, Title, Content, Created_At FROM posts
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func handlePostPost(w http.ResponseWriter, r *http.Request) {
    var post Post
    err := json.NewDecoder(r.Body).Decode(&post)
    if err != nil {
        http.Error(w, `{"error":"invalid JSON: `+err.Error()+`"}`, http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(post)
}

func handleGetPostByID(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    row := db.QueryRow(`
        SELECT id, title, content, user_id, created_at
        FROM posts
        WHERE id = ?
    `, id)

    var post Post
    err = row.Scan(
        &post.ID,
        &post.Title,
        &post.Content,
        &post.UserID,
        &post.CreatedAt,
    )

    if err == sql.ErrNoRows {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}

//COMMENTS
func commentHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        handleGetComment(w, r)
    case http.MethodPost:
        handlePostComment(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func handleGetComment(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
        SELECT ID, Post_ID, User_ID, Comment, Created_At FROM comments
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var comment Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Comment,
			&comment.CreatedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func handlePostComment(w http.ResponseWriter, r *http.Request) {
    var comment Comment
    err := json.NewDecoder(r.Body).Decode(&comment)
    if err != nil {
        http.Error(w, `{"error":"invalid JSON: `+err.Error()+`"}`, http.StatusBadRequest)
        return
    }

    json.NewEncoder(w).Encode(comment)
}

func handleGetCommentByID(w http.ResponseWriter, r *http.Request) {
    idStr := strings.TrimPrefix(r.URL.Path, "/comments/")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    row := db.QueryRow(`
        SELECT id, post_id, user_id, comment, created_at
        FROM comments
        WHERE id = ?
    `, id)

    var comment Comment
    err = row.Scan(
        &comment.ID,
        &comment.PostID,
        &comment.UserID,
        &comment.Comment,
        &comment.CreatedAt,
    )

    if err == sql.ErrNoRows {
        http.Error(w, "Comment not found", http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comment)
}

