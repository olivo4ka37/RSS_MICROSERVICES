package auth

import (
	"api-service/app/db"
	"context"
	"github.com/google/uuid"
	"log"
)

func UserAuth(uuidStr string) (bool, error) {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		log.Printf("Invalid UUID: %v", err)
		return false, err
	}

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE uuid = $1)"
	err = db.Conn.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return false, err
	}

	return exists, nil
}

func AdminAuth(uuidStr string) (bool, error) {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		log.Printf("Invalid UUID: %v", err)
		return false, err
	}

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM administrators WHERE uuid = $1)"
	err = db.Conn.QueryRow(context.Background(), query, id).Scan(&exists)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return false, err
	}

	return exists, nil
}
