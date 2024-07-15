package handlers

import (
	"api-service/app/auth"
	"api-service/app/db"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GetAllSources(w http.ResponseWriter, r *http.Request) {
	authID := r.Header.Get("X-Auth-ID")
	if authID == "" {
		http.Error(w, "Missing X-Auth-ID header", http.StatusUnauthorized)
		return
	}

	authenticated, err := auth.UserAuth(authID)
	if err != nil {
		log.Println("Error while trying to Authenticate user:", err)
		return
	}

	if !authenticated {

		authenticated, err = auth.AdminAuth(authID)
		if err != nil {
			log.Println("Error while trying to Authenticate admin:", err)
			return
		}

		if !authenticated {
			http.Error(w, "User with uuid is not existr", http.StatusUnauthorized)
			return
		}
	}

	url := "http://localhost:8080/sources"
	resp, err := http.Get(cmp.Or(os.Getenv("srcURL"), url))
	if err != nil {
		http.Error(w, "Failed to get all sources:", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var apiResp []db.Source
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("failed to decode response: %w", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apiResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//w.WriteHeader(resp.StatusCode)
}

// GetUserNews fetches all news from subscribed sources of a user
func GetUserNews(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Header.Get("X-Auth-ID")
	if userUUID == "" {
		http.Error(w, "Missing X-Auth-ID header", http.StatusUnauthorized)
		return
	}

	var id int

	err := db.Conn.QueryRow(context.Background(), "SELECT id FROM users WHERE uuid=$1", userUUID).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot find user with this uuid (Query failed): %v\n", err)
		os.Exit(1)
	}

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

	if userID != id {
		http.Error(w, "You trying to view subscriptions of another person", http.StatusBadRequest)
		return
	}

	// Get pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Set default values if parameters are not provided
	page := 1
	limit := 2

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	url := fmt.Sprintf("http://news-service:8080/sources/%d/news?page=%d&limit=%d", id, page, limit)
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to get all subscribed sources: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var apiResp []db.Article
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("failed to decode response: %w", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(apiResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//w.WriteHeader(resp.StatusCode)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	userUUID := uuid.NewString()
	var userID int
	err := db.Conn.QueryRow(context.Background(), "INSERT INTO users (uuid) VALUES ($1) RETURNING id", userUUID).Scan(&userID)
	if err != nil {
		fmt.Errorf("error inserting user: %v", err)
		return
	}

	log.Println(auth.UserAuth(userUUID))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userUUID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AddSource(w http.ResponseWriter, r *http.Request) {

	authID := r.Header.Get("X-Auth-ID")
	if authID == "" {
		http.Error(w, "Missing X-Auth-ID header", http.StatusUnauthorized)
		return
	}

	authenticated, err := auth.AdminAuth(authID)
	if err != nil {
		log.Println("Error while trying to Authenticate admin:", err)
		return
	}
	if !authenticated {
		http.Error(w, "User with uuid is not existr", http.StatusUnauthorized)
		return
	}

	url := "http://localhost:8080/sources"
	resp, err := http.Post(cmp.Or(os.Getenv("srcURL"), url), "application/json", r.Body)
	if err != nil {
		http.Error(w, "Failed to Add source:", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
}

func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	adminUUID := uuid.NewString()
	var userID int
	err := db.Conn.QueryRow(context.Background(), "INSERT INTO administrators (uuid) VALUES ($1) RETURNING id", adminUUID).Scan(&userID)
	if err != nil {
		fmt.Errorf("error inserting user: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(adminUUID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Subscribe a user to an RSS source
func SubscribeUser(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Header.Get("X-Auth-ID")
	if userUUID == "" {
		http.Error(w, "Missing X-Auth-ID header", http.StatusUnauthorized)
		return
	}

	var userID int
	err := db.Conn.QueryRow(context.Background(), "SELECT id FROM users WHERE uuid = $1", userUUID).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

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

	var sourceID int
	err = db.Conn.QueryRow(context.Background(), "SELECT id FROM sources WHERE url = $1", source.URL).Scan(&sourceID)
	if err != nil {
		http.Error(w, "Source not found", http.StatusNotFound)
		return
	}

	_, err = db.Conn.Exec(context.Background(), "INSERT INTO user_subscriptions (user_id, source_id) VALUES ($1, $2)", userID, sourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Unsubscribe a user from an RSS source
func UnsubscribeUser(w http.ResponseWriter, r *http.Request) {
	userUUID := r.Header.Get("X-Auth-ID")
	if userUUID == "" {
		http.Error(w, "Missing X-Auth-ID header", http.StatusUnauthorized)
		return
	}

	var userID int
	err := db.Conn.QueryRow(context.Background(), "SELECT id FROM users WHERE uuid = $1", userUUID).Scan(&userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

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

	var sourceID int
	err = db.Conn.QueryRow(context.Background(), "SELECT id FROM sources WHERE url = $1", source.URL).Scan(&sourceID)
	if err != nil {
		http.Error(w, "Source not found", http.StatusNotFound)
		return
	}

	_, err = db.Conn.Exec(context.Background(), "DELETE FROM user_subscriptions WHERE user_id=$1 AND source_id=$2", userID, sourceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
