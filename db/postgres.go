package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		connStr = "postgres://divyam.sinha@localhost:5432/postgres?sslmode=disable"
	}

	var db *sql.DB
	var err error

	for i := 0; i < 15; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			if err = db.Ping(); err == nil {
				log.Println("Connected to Postgres")
				return db
			}
		}

		log.Printf("Failed to connect to postgres: %v. Retrying in 2 seconds (%d/15)...", err, i+1)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Could not connect to postgres after multiple retries: ", err)
	return nil
}
