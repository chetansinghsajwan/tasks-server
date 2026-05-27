package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init() {

	var dbUser = os.Getenv("DB_USER")
	if dbUser == "" {
		panic("DB_USER environment variable is not set")
	}

	var dbPassword = os.Getenv("DB_PASS")
	if dbPassword == "" {
		panic("DB_PASS environment variable is not set")
	}

	var dbName = os.Getenv("DB_NAME")
	if dbName == "" {
		panic("DB_NAME environment variable is not set")
	}

	var dbHost = os.Getenv("DB_HOST")
	if dbHost == "" {
		panic("DB_HOST environment variable is not set")
	}

	var dbPort = os.Getenv("DB_PORT")
	if dbPort == "" {
		panic("DB_PORT environment variable is not set")
	}

	var dbSSLMode = os.Getenv("DB_SSL_MODE")
	if dbSSLMode == "" {
		panic("DB_SSL_MODE environment variable is not set")
	}

	var connString = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbSSLMode,
	)

	print(connString)

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
