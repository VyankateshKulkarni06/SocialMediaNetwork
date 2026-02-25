package db

import (
	"database/sql"
	"fmt"
	"log"
)

func InitSchema(database *sql.DB) {

	createUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		bio TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	createConnections := `
	CREATE TABLE IF NOT EXISTS connections (
		id BIGSERIAL PRIMARY KEY,
		from_user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
		to_user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT NOW(),
		UNIQUE (from_user_id, to_user_id)
	);`

	_, err := database.Exec(createUsers)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.Exec(createConnections)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Schema initialized")
}