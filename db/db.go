package db

import (
	"context"
	"os"
	"tasks/sqlc"

	// "github.com/jackc/pgx"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Queries *sqlc.Queries
var Pool *pgxpool.Pool

func Init() {

	var connString = os.Getenv("DB_CONN_STRING")
	if connString == "" {
		panic("DB_CONN_STRING environment variable is not set")
	}

	var ctx = context.Background()

	var err error
	Pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	if err := Pool.Ping(ctx); err != nil {
		panic("Failed to ping the database: " + err.Error())
	}

	Queries = sqlc.New(Pool)

	println("Database connection established successfully")
}

type TxQueries struct {
	Queries *sqlc.Queries
	Tx      pgx.Tx
}

func Begin(ctx context.Context) (*TxQueries, error) {

	var tx, err = Pool.Begin(ctx)

	if err != nil {
		return nil, err
	}

	var query = TxQueries{
		Queries: sqlc.New(tx),
		Tx:      tx,
	}

	return &query, nil
}
