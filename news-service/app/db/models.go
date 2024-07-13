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

var startSources []Source = []Source{
	{URL: "https://habr.com/ru/rss/hub/go/all/?fl=ru"},
	{URL: "https://habr.com/ru/rss/best/daily/?fl=ru"},
	{URL: "https://golangcode.com/index.xml"},
	{URL: "https://forum.golangbridge.org/latest.rss"},
	{URL: "https://appliedgo.net/index.xml"},
	{URL: "https://blog.jetbrains.com/go/feed/"},
	{URL: "https://dave.cheney.net/category/golang/feed"},
	{URL: "https://changelog.com/gotime/feed"},
	{URL: "https://golang.ch/feed/"},
	{URL: "https://gosamples.dev/index.xml"},
}
