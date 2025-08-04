package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init() {

	var connString = os.Getenv("DATABASE_URL")
	if connString == "" {
		panic("DATABASE_URL environment variable is not set")
	}

	var ctx = context.Background()

	var err error
	Pool, err = pgxpool.New(ctx, connString)
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	if err = Pool.Ping(ctx); err != nil {
		panic("Failed to ping the database: " + err.Error())
	}

	println("Database connection established successfully")
}
