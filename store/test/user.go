package test

import (
	"context"
	"tasks/errorcodes"
	"tasks/option"
	"tasks/store"
	"testing"
)

func TestUserStore(t *testing.T, st store.Store) {

	var User0_ID = "user0"
	var User0_Email = "user0@domain.com"
	var User0_FullName = "First0 Middle0 Last0"
	var User0_DisplayName = "First0"

	var User1_ID = "user1"
	var User1_Email = "user1@domain.com"
	var User1_FullName = "First1 Middle1 Last1"
	var User1_DisplayName = "First1"

	var ctx = context.Background()

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid id
	// ---------------------------------------------------------------------------------------------

	var serr *store.StoreError
	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "",
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid id
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "   ",
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid id
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "-username",
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid id
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "username-",
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid id
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "user--name",
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid id
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          " username",
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid id
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "username ",
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserIDFormat {
		t.Fatalf("Expected store.ErrUserIDFormatCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid full name
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "username0",
		Email:       User0_Email,
		FullName:    "  ",
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserFullNameFormat {
		t.Fatalf("Expected store.ErrUserFullNameFormatCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid full name
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "username0",
		Email:       User0_Email,
		FullName:    User0_FullName + " ",
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserFullNameFormat {
		t.Fatalf("Expected store.ErrUserFullNameFormatCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid full name
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "username0",
		Email:       User0_Email,
		FullName:    " " + User0_FullName,
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserFullNameFormat {
		t.Fatalf("Expected store.ErrUserFullNameFormatCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid display name
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "username0",
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(" "),
	})

	if serr == nil || serr.Code != errorcodes.UserDisplayNameFormat {
		t.Fatalf("Expected errorcodes.UserDisplayNameFormat error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid display name
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "username0",
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(User0_DisplayName + " "),
	})

	if serr == nil || serr.Code != errorcodes.UserDisplayNameFormat {
		t.Fatalf("Expected errorcodes.UserDisplayNameFormat error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with invalid display name
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          "username0",
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(" " + User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserDisplayNameFormat {
		t.Fatalf("Expected errorcodes.UserDisplayNameFormat error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          User0_ID,
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr != nil {
		t.Fatalf("Expected nil err, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test get user that doesn't exist
	// ---------------------------------------------------------------------------------------------

	_, serr = st.GetUser(ctx, "usernamenotexists")

	if serr == nil || serr.Code != errorcodes.UserNotFound {
		t.Fatalf("Expected errorcodes.UserNotFound error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test create user with same id again
	// ---------------------------------------------------------------------------------------------

	serr = st.CreateUser(ctx, store.CreateUserParams{
		ID:          User0_ID,
		Email:       User0_Email,
		FullName:    User0_FullName,
		DisplayName: option.Some(User0_DisplayName),
	})

	if serr == nil || serr.Code != errorcodes.UserIDAlreadyExists {
		t.Fatalf("Expected store.ErrUserIDAlreadyExistsCode error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Test get user
	// ---------------------------------------------------------------------------------------------

	var user *store.User
	user, serr = st.GetUser(ctx, User0_ID)

	if serr != nil {
		t.Fatalf("Expected nil err, got: %v", serr)
	}

	if user.Email != User0_Email {
		t.Fatalf("Expected user.email '%s', got: %s", User0_Email, user.Email)
	}

	if user.FullName != User0_FullName {
		t.Fatalf("Expected user.email '%s', got: %s", User0_FullName, user.FullName)
	}

	if user.DisplayName.MustGet() != User0_DisplayName {
		t.Fatalf("Expected user.email '%s', got: %s", User0_DisplayName, user.DisplayName.MustGet())
	}

	// ---------------------------------------------------------------------------------------------
	// Update user with invalid FullName
	// ---------------------------------------------------------------------------------------------

	var fullName = " "
	serr = st.UpdateUser(ctx, User0_ID, store.UpdateUserParams{
		FullName: option.Some(&fullName),
	})

	if serr == nil || serr.Code != errorcodes.UserFullNameFormat {
		t.Fatalf("Expected errorcodes.UserFullNameFormat error, got: %v", serr)
	}

	// ---------------------------------------------------------------------------------------------
	// Update user
	// ---------------------------------------------------------------------------------------------

	serr = st.UpdateUser(ctx, User0_ID, store.UpdateUserParams{
		ID:          option.Some(&User1_ID),
		Email:       option.Some(&User1_Email),
		FullName:    option.Some(&User1_FullName),
		DisplayName: option.Some(&User1_DisplayName),
	})

	if serr != nil {
		t.Fatalf("Expected nil error, got: %v", serr.WrappedError.Error())
	}

	// ---------------------------------------------------------------------------------------------
	// Check if the user values are correct
	// ---------------------------------------------------------------------------------------------

	user, serr = st.GetUser(ctx, User1_ID)

	if serr != nil {
		t.Fatalf("Expected nil err, got: %v", serr)
	}

	if user.ID != User1_ID {
		t.Fatalf("Expected User.ID '%s', got: '%s'", User1_ID, user.ID)
	}

	if user.Email != User1_Email {
		t.Fatalf("Expected User.Email '%s', got: '%s'", User1_Email, user.Email)
	}

	if user.FullName != User1_FullName {
		t.Fatalf("Expected User.FullName '%s', got: '%s'", User1_FullName, user.FullName)
	}

	if user.DisplayName.MustGet() != User1_DisplayName {
		t.Fatalf("Expected User.DisplayName '%s', got: '%s'", User1_DisplayName, user.DisplayName.MustGet())
	}

	// ---------------------------------------------------------------------------------------------
	// Test delete user
	// ---------------------------------------------------------------------------------------------

	serr = st.DeleteUser(ctx, User1_ID)

	if serr != nil {
		t.Fatalf("Expected nil err, got: %v", serr.WrappedError.Error())
	}

	// ---------------------------------------------------------------------------------------------
	// Test delete user that doesn't exist
	// ---------------------------------------------------------------------------------------------

	serr = st.DeleteUser(ctx, User1_ID)

	if serr == nil || serr.Code != errorcodes.UserNotFound {
		t.Fatalf("Expected errorcodes.UserNotFound error, got: %v", serr)
	}
}
