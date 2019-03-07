package common

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(databaseUrl string) *sql.DB {
	db, err := sql.Open("postgres", databaseUrl)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func ClearDataTestDB(db *sql.DB) {
	_, err := db.Exec(`TRUNCATE users`)

	if err != nil {
		panic(err)
	}
}
