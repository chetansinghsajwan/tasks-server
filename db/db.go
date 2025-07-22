package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {

	var connString = os.Getenv("DB_CONN_STRING")
	if connString == "" {
		panic("DB_CONN_STRING environment variable is not set")
	}

	var error error
	DB, error = sql.Open("postgres", connString)
	if error != nil {
		panic("Failed to connect to the database: " + error.Error())
	}

	println("Database connection established successfully")
}
