package pg_test

import (
	"context"
	"tasks/option"
	"tasks/store"
	"tasks/store/pg"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCreateUser(t *testing.T) {

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
	var us = pg.PgUserStore{Pool: pool}

	// Test create user
	var serr *store.StoreError
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "",
		Email:       "myemail@domain.com",
		FullName:    "First Middle Last",
		DisplayName: option.Some("First"),
	})

	if serr.Code != store.ErrUserIDFormatCode {
		t.Fatalf("Expected store.ErrUserIDFormat error, got: %v", serr)
	}

	// Test create user
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "   ",
		Email:       "myemail@domain.com",
		FullName:    "First Middle Last",
		DisplayName: option.Some("First"),
	})

	if serr.Code != store.ErrUserIDFormatCode {
		t.Fatalf("Expected store.ErrUserIDFormat error, got: %v", serr)
	}

	// Test create user
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "-username",
		Email:       "myemail@domain.com",
		FullName:    "First Middle Last",
		DisplayName: option.Some("First"),
	})

	if serr.Code != store.ErrUserIDFormatCode {
		t.Fatalf("Expected store.ErrUserIDFormat error, got: %v", serr)
	}

	// Test create user
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "username-",
		Email:       "myemail@domain.com",
		FullName:    "First Middle Last",
		DisplayName: option.Some("First"),
	})

	if serr.Code != store.ErrUserIDFormatCode {
		t.Fatalf("Expected store.ErrUserIDFormat error, got: %v", serr)
	}

	// Test create user
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "user--name",
		Email:       "myemail@domain.com",
		FullName:    "First Middle Last",
		DisplayName: option.Some("First"),
	})

	if serr.Code != store.ErrUserIDFormatCode {
		t.Fatalf("Expected store.ErrUserIDFormat error, got: %v", serr)
	}

	// Test create user
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "username",
		Email:       "myemail@domain.com",
		FullName:    "First Middle Last",
		DisplayName: option.Some("First"),
	})

	if serr != nil {
		t.Fatalf("Expected nil err, got: %v", serr)
	}
}
