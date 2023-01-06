package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func StartDB() {
	var err error
	url := os.Getenv("DATABASE_URL")
	DB, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal("connection to database error", err)
		return
	}

	createTable := `CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);`

	_, err = DB.Exec(createTable)

	if err != nil {
		log.Fatal("connection to database error", err)
		return
	}
}
