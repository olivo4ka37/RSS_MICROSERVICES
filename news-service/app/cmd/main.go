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
