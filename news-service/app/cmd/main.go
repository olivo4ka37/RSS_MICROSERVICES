package main

import (
	"context"
	_ "github.com/lib/pq"
	"log"
	"news_service/app/cache"
	"news_service/app/db"
	"news_service/app/server"
	"news_service/app/workers"
)

/*
type Source struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type RSSItem struct {
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
	{URL: "http://sdet.us/category/golang/feed/"},
	{URL: "https://gosamples.dev/index.xml"},
}
*/

var rssCache = cache.RssCache{}

func main() {
	var err error
	db.Conn, err = db.ConnAndLoad()
	if err != nil {
		log.Fatalf("Error while trying to connect to db:", err)
	}
	defer db.Conn.Close(context.Background())

	go workers.StartWorkers(&rssCache)

	log.Fatalf("server stopped work:", server.NewServer())
}
