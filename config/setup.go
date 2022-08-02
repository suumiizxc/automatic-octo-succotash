package config

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

const (
	host     = "0.0.0.0"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "nextcloud"
)

func OpenConnection() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	DB = db
}
