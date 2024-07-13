package db

import "time"

type Source struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type User struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	LastLogin time.Time `json:"last_login"`
}

type RSSItem struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Published   time.Time `json:"published"`
	SourceID    int       `json:"source_id"`
}

type Article struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Published   time.Time `json:"published"`
	SourceID    int       `json:"source_id"`
}
