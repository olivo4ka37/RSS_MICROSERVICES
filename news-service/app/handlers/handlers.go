package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"news_service/app/db"
	"strconv"
)

func GetSourcesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Conn.Query(context.Background(), "SELECT id, url FROM sources")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sources []db.Source
	for rows.Next() {
		var source db.Source
		if err := rows.Scan(&source.ID, &source.URL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sources = append(sources, source)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sources); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetUserNews(w http.ResponseWriter, r *http.Request) {
	// Extract userID from URL
	vars := mux.Vars(r)
	userIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get pagination parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		http.Error(w, "Error while getting number of page", http.StatusBadRequest)
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "Error while getting number of limit", http.StatusBadRequest)
	}
	log.Println("page is:", page)
	log.Println("limit is:", limit)

	// Calculate offset
	offset := (page - 1) * limit

	// Get articles for the sources after last login time with pagination
	query := `SELECT a.id, a.title, a.link, a.description, a.published, a.source_id
FROM articles a
JOIN user_subscriptions us ON a.source_id = us.source_id
JOIN users u ON us.user_id = u.id
WHERE us.user_id = $1
  AND a.published > u.last_login
ORDER BY a.published DESC
LIMIT $2 OFFSET $3;`
	articleRows, err := db.Conn.Query(context.Background(), query, userID, limit, offset)
	if err != nil {
		log.Println("Error while trying to get articleRows")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Article rows is: ", articleRows)
	defer articleRows.Close()

	articles := make([]db.Article, 0, 1000)
	for articleRows.Next() {
		var article db.Article
		if err := articleRows.Scan(&article.ID, &article.Title, &article.Link, &article.Description, &article.Published, &article.SourceID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		articles = append(articles, article)
	}
	log.Println("articles is:", articles)

	// Return articles as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func AddSourceHandler(w http.ResponseWriter, r *http.Request) {
	var source db.Source
	if err := json.NewDecoder(r.Body).Decode(&source); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(source.URL) == 0 {
		log.Println("No content")
		w.WriteHeader(204)
		return
	}

	err := db.Conn.QueryRow(context.Background(), "INSERT INTO sources (url) VALUES ($1) RETURNING id", source.URL).Scan(&source.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
