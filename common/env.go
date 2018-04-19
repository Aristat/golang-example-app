package common

import (
	"database/sql"
)

type Env struct {
	DB *sql.DB
}

func InitEnv() *Env {
	env := &Env{DB: initDB()}
	return env
}

func InitTestEnv() *Env {
	env := &Env{DB: initTestDB()}
	return env
}
