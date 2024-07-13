package main

import (
	"api-service/app/db"
	"api-service/app/server"
	"context"
	"log"
)

func main() {
	var err error
	db.Conn, err = db.ConnToDB()
	if err != nil {
		log.Fatalf("Error while trying to connect to db:", err)
	}
	defer db.Conn.Close(context.Background())

	log.Fatalf("server stopped work:", server.NewServer())
}
