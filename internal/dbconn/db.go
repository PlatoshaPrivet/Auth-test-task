package dbconn

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func Open() *sql.DB {
	dbURL := os.Getenv("DATABASE_URL") //taking URL from docker-compose

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	return db
}
func CreateTable(db *sql.DB) {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS refresh_tokens (
        user_id UUID NOT NULL UNIQUE,
        email TEXT NOT NULL UNIQUE,
        refresh_token TEXT NOT NULL
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
}

//dbURL := "postgres://postgres:onlyfortesting52528249@localhost:5432/test?sslmode=disable"
