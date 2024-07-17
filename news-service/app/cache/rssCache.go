package cache

import (
	"context"
	"log"
	"news_service/app/db"
	"sync"
	"time"
)

type RssCache struct {
	Sources []db.Source
	Mut     sync.RWMutex
}

func (rc *RssCache) CacheWorker() {
	time.Sleep(time.Second * 5)
	rc.update()
}

func (rc *RssCache) reset() {
	rc = &RssCache{}
}

func (rc *RssCache) store() {
	rows, err := db.Conn.Query(context.Background(), "SELECT id, url FROM Sources")
	if err != nil {
		log.Printf("Error querying Sources: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var source db.Source
		if err := rows.Scan(&source.ID, &source.URL); err != nil {
			log.Printf("Error scanning source: %v\n", err)
			continue
		}

		rc.Sources = append(rc.Sources, source)
	}
}

func (rc *RssCache) update() {
	rc.reset()
	rc.Mut.Lock()
	rc.store()
	rc.Mut.Unlock()

	log.Println("Cache updated")
}
