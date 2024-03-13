package user

import (
	"database/sql"
	"log"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=user sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
}
