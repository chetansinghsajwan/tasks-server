package pg_test

import (
	"context"
	"tasks/option"
	"tasks/store"
	"tasks/store/pg"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const User1_ID = "username"
const User1_Email = "myemail@domain.com"
const User1_FullName = "First Middle Last"
const User1_DisplayName = "First"

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
	var us = pg.PostgresStore{Pool: pool}

	// Test create user with invalid id
	var serr *store.StoreError
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "",
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// Test create user with invalid id
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "   ",
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// Test create user with invalid id
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "-username",
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// Test create user with invalid id
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "username-",
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// Test create user with invalid id
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "user--name",
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// Test create user with invalid id
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          " username",
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// Test create user with invalid id
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "username ",
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// Test create user with invalid full name
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "username1",
		Email:       User1_Email,
		FullName:    "  ",
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserFullNameFormat {
		t.Fatalf("Expected store.ErrUserFullNameFormatCode error, got: %v", serr)
	}

	// Test create user with invalid full name
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "username1",
		Email:       User1_Email,
		FullName:    User1_FullName + " ",
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserFullNameFormat {
		t.Fatalf("Expected store.ErrUserFullNameFormatCode error, got: %v", serr)
	}

	// Test create user with invalid full name
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "username1",
		Email:       User1_Email,
		FullName:    " " + User1_FullName,
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserFullNameFormat {
		t.Fatalf("Expected store.ErrUserFullNameFormatCode error, got: %v", serr)
	}

	// Test create user with invalid display name
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "username1",
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(" "),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserDisplayNameFormat {
		t.Fatalf("Expected store.ErrorCode_UserDisplayNameFormat error, got: %v", serr)
	}

	// Test create user with invalid display name
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "username1",
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(User1_DisplayName + " "),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserDisplayNameFormat {
		t.Fatalf("Expected store.ErrorCode_UserDisplayNameFormat error, got: %v", serr)
	}

	// Test create user with invalid display name
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          "username1",
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(" " + User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserDisplayNameFormat {
		t.Fatalf("Expected store.ErrorCode_UserDisplayNameFormat error, got: %v", serr)
	}

	// Test create user
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          User1_ID,
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr != nil {
		t.Fatalf("Expected nil err, got: %v", serr)
	}

	// Test get user that doesn't exist
	_, serr = us.GetUser(ctx, "usernamenotexists")

	if serr == nil || serr.Code != store.ErrorCode_UserNotFound {
		t.Fatalf("Expected store.ErrorCode_UserNotFound error, got: %v", serr)
	}

	// Test create user with same id again
	serr = us.CreateUser(ctx, store.CreateUserParams{
		ID:          User1_ID,
		Email:       User1_Email,
		FullName:    User1_FullName,
		DisplayName: option.Some(User1_DisplayName),
	})

	if serr == nil || serr.Code != store.ErrorCode_UserIDAlreadyExists {
		t.Fatalf("Expected store.ErrUserIDAlreadyExistsCode error, got: %v", serr)
	}

	// Test get user
	var user *store.User
	user, serr = us.GetUser(ctx, User1_ID)

	if serr != nil {
		t.Fatalf("Expected nil err, got: %v", serr)
	}

	if user.Email != User1_Email {
		t.Fatalf("Expected user.email '%s', got: %s", User1_Email, user.Email)
	}

	if user.FullName != User1_FullName {
		t.Fatalf("Expected user.email '%s', got: %s", User1_FullName, user.FullName)
	}

	if user.DisplayName.MustGet() != User1_DisplayName {
		t.Fatalf("Expected user.email '%s', got: %s", User1_DisplayName, user.DisplayName.MustGet())
	}

	// Test delete user
	serr = us.DeleteUser(ctx, User1_ID)

	if serr != nil {
		t.Fatalf("Expected nil err, got: %v", serr.WrappedError.Error())
	}

	// Test delete user that doesn't exist
	serr = us.DeleteUser(ctx, User1_ID)

	if serr == nil || serr.Code != store.ErrorCode_UserNotFound {
		t.Fatalf("Expected store.ErrorCode_UserNotFound error, got: %v", serr)
	}
}
