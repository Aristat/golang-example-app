package common

import (
	"database/sql"

	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func initDB() *sql.DB {
	db, err := sql.Open("postgres", viper.GetString("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func initTestDB() *sql.DB {
	db, err := sql.Open("postgres", viper.GetString("TEST_DATABASE_URL"))

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
