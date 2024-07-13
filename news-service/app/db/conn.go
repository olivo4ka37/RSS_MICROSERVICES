package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"time"
)

var counts int
var Conn *pgx.Conn

func openDB(dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func connToDB() (*pgx.Conn, error) {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to postgres (YOPTA)!")
			return connection, nil
		}

		if counts > 10 {
			log.Println(err)
			return nil, err
		}

		log.Println("Backing off (Chill) 2 seconds")
		time.Sleep(time.Second * 2)
		continue
	}

	/*
		conn, err := pgx.Connect(context.Background(), config.ConnStr)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
			return nil, err
		}

		if err = loadStartingSources(conn); err != nil{
			log.Fatalf("Error while trying to load start RSS sources: %v\n", err)
			return nil, err
		}

		log.Println("Connected to DB and loaded starting sources!")

	*/

	//return conn, nil
}

func ConnAndLoad() (*pgx.Conn, error) {
	conn, err := connToDB()

	if err = loadStartingSources(conn); err != nil {
		log.Fatalf("Error while trying to load start RSS sources: %v\n", err)
		return nil, err
	}

	log.Println("Connected to DB and loaded starting sources!")

	return conn, nil
}

func loadStartingSources(conn *pgx.Conn) error {
	for _, src := range startSources {
		err := conn.QueryRow(context.Background(), `
            INSERT INTO sources (url) 
            VALUES ($1) 
            ON CONFLICT (url) DO NOTHING 
            RETURNING id`, src.URL).Scan(&src.ID)

		if err != nil {
			if err == pgx.ErrNoRows {
				// Получаем id источника из базы данных, если он уже существует
				err = conn.QueryRow(context.Background(), "SELECT id FROM sources WHERE url=$1", src.URL).Scan(&src.ID)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		//fetchAndStoreRSS(src)
	}

	return nil
}
