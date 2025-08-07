package pg_test

import (
	"context"
	"tasks/store"
	"tasks/store/pg"
	"tasks/store/test"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresUserSecretStore(t *testing.T) {

	var ctx = context.Background()
	var connString = "postgres://devuser:devpass@database:5432/testdb?sslmode=disable"

	// Setup the connection
	var pool *pgxpool.Pool
	var err error
	pool, err = pgxpool.New(ctx, connString)

	if err != nil {
		t.Fatalf("Failed to connect to the database, error: %s", err.Error())
	}

	// Setup the database
	_, err = pool.Exec(ctx, "select reinitialize_schema()")
	if err != nil {
		t.Fatalf("Failed to reinitialize the database, error: %s", err.Error())
	}

	// Setup stores
	var st store.Store = pg.PostgresStore{Pool: pool}

	// Perform the test
	test.TestUserSecretStore(t, st)
}
