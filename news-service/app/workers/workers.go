package workers

import (
	"context"
	"github.com/mmcdole/gofeed"
	"log"
	"news_service/app/cache"
	"news_service/app/db"
	"time"
)

func StartWorkers(rc *cache.RssCache) {
	for {
		rc.CacheWorker()
		updateItemsWorker(rc)
	}
}

func updateItemsWorker(rc *cache.RssCache) {
	time.Sleep(time.Second * 5)
	rc.Mut.RLock()
	for _, src := range rc.Sources {
		getItemsFromRSS(src)
	}
	rc.Mut.RUnlock()
	log.Println("All items are updated!")
}

func getItemsFromRSS(src db.Source) {
	// fp - feed parser
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(src.URL)
	if err != nil {
		log.Printf("Error fetching RSS feed from %s: %v", src.URL, err)
		return
	}

	for i, item := range feed.Items {
		if i >= 2 {
			break
		}

		/*
			rows, err := db.Conn.Query(context.Background(), "INSERT INTO articles (title, link, description, published, source_id) VALUES ($1, $2, $3, $4, $5)",
				item.Title, item.Link, item.Description, item.PublishedParsed, src.ID)
			if err != nil {
				log.Printf("Error inserting article: %v", err)
				return
			}
			rows.Close()
		*/

		_, err := db.Conn.Exec(context.Background(), `
        INSERT INTO articles (title, link, description, published, source_id) 
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (link) DO UPDATE 
        SET title = EXCLUDED.title, 
            description = EXCLUDED.description, 
            published = EXCLUDED.published, 
            source_id = EXCLUDED.source_id`,
			item.Title, item.Link, item.Description, item.PublishedParsed, src.ID)
		if err != nil {
			log.Printf("Error inserting or updating article: %v", err)
			return
		}
	}
}
