package pg

import "github.com/jackc/pgx/v5/pgxpool"

type PostgresStore struct {
	Pool *pgxpool.Pool
}
