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
}

func ConnAndLoad() (*pgx.Conn, error) {
	conn, err := connToDB()
	if err != nil {
		log.Fatalf("Can't connect to DB: %e", err)
	}

	log.Println("Connected to DB and loaded starting sources!")

	return conn, nil
}
