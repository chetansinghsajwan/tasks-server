package db

import (
	"context"
	"os"
	"tasks/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Queries *sqlc.Queries

func Init() {

	var connString = os.Getenv("DB_CONN_STRING")
	if connString == "" {
		panic("DB_CONN_STRING environment variable is not set")
	}

	var ctx = context.Background()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	if err := pool.Ping(ctx); err != nil {
		panic("Failed to ping the database: " + err.Error())
	}

	Queries = sqlc.New(pool)

	println("Database connection established successfully")
}
