package db

import (
	"cmp"
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
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func ConnToDB() (*pgx.Conn, error) {
	dsn := cmp.Or(os.Getenv("DSN"), "host=localhost port=5040 user=postgres password=password dbname=rss sslmode=disable timezone=UTC connect_timeout=5")

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
