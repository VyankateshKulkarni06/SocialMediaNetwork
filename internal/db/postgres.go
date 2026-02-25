package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName,
	)

	var database *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		database, err = sql.Open("postgres", connStr)
		if err == nil {
			err = database.Ping()
			if err == nil {
				break
			}
		}

		fmt.Println("Waiting for DB...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	fmt.Println("Connected to PostgreSQL")
	return database
}