package test

import (
	"context"
	"tasks/option"
	"tasks/store"
	"testing"
)

func TestUserSecretStore(t *testing.T, st store.Store) {

	var ValidUserID_0, _ = store.ParseUserID("myid_0")
	var ValidSecretPass_0 = "mypassword_0"

	var ValidUserID_1, _ = store.ParseUserID("myid_1")
	var ValidSecretPass_1 = "mypassword_1"

	var ctx = context.Background()

	st.CreateUser(ctx, store.CreateUserParams{
		ID:          ValidUserID_0,
		Email:       "email@domain.com",
		FullName:    "First Middle Last",
		DisplayName: option.Some("First"),
	})

	t.Run("Create secret non existing user", func(t *testing.T) {

		var serr *store.StoreError = st.CreateUserSecret(ctx, store.CreateUserSecretParams{
			ID:   ValidUserID_1,
			Pass: ValidSecretPass_0,
		})

		if serr == nil || serr.Code != store.ErrorCode_UserNotFound {
			t.Fatalf("Expected store.ErrorCode_UserNotFound error, got %v", serr)
		}
	})

	t.Run("Create secret with invalid pass", func(t *testing.T) {

		var serr *store.StoreError = st.CreateUserSecret(ctx, store.CreateUserSecretParams{
			ID:   ValidUserID_0,
			Pass: "  ",
		})

		if serr == nil || serr.Code != store.ErrorCode_InvalidUserSecretPassFormat {
			t.Fatalf("Expected store.ErrorCode_InvalidSecretValueFormat error, got %v", serr)
		}
	})

	t.Run("Create secret with valid fields", func(t *testing.T) {

		var serr *store.StoreError = st.CreateUserSecret(ctx, store.CreateUserSecretParams{
			ID:   ValidUserID_0,
			Pass: ValidSecretPass_0,
		})

		if serr != nil {
			t.Fatalf("Expected nil error, got %v", serr)
		}
	})

	t.Run("Get secret that doesn't exist", func(t *testing.T) {

		var serr *store.StoreError

		_, serr = st.GetUserSecret(ctx, ValidUserID_1)

		if serr == nil || serr.Code != store.ErrorCode_UserSecretNotFound {
			t.Fatalf("Expected store.ErrorCode_SecretNotFound error, got %v", serr)
		}
	})

	t.Run("Get valid secret", func(t *testing.T) {

		var secret *store.UserSecret
		var serr *store.StoreError

		secret, serr = st.GetUserSecret(ctx, ValidUserID_0)

		if serr != nil {
			t.Fatalf("Expected nil error, got %v", serr)
		}

		if secret.ID != ValidUserID_0 {
			t.Fatalf("Expected secret.ID '%s', got '%s'", ValidUserID_0, secret.ID)
		}

		if secret.Pass != ValidSecretPass_0 {
			t.Fatalf("Expected secret.Pass '%s', got '%s'", ValidSecretPass_0, secret.Pass)
		}
	})

	t.Run("Update secret for non existing user", func(t *testing.T) {

		var serr *store.StoreError = st.UpdateUserSecret(ctx, ValidUserID_1,
			store.UpdateUserSecretParams{
				Pass: ValidSecretPass_1,
			},
		)

		if serr == nil || serr.Code != store.ErrorCode_UserNotFound {
			t.Fatalf("Expected store.ErrorCode_UserNotFound error, got %v", serr)
		}
	})

	t.Run("Update secret with invalid pass", func(t *testing.T) {

		var serr *store.StoreError = st.UpdateUserSecret(ctx, ValidUserID_0,
			store.UpdateUserSecretParams{
				Pass: "  ",
			},
		)

		if serr == nil || serr.Code != store.ErrorCode_InvalidUserSecretPassFormat {
			t.Fatalf("Expected store.ErrorCode_InvalidSecretValueFormat error, got %v", serr)
		}
	})

	t.Run("Update secret with valid fields", func(t *testing.T) {

		var serr *store.StoreError = st.UpdateUserSecret(ctx, ValidUserID_0,
			store.UpdateUserSecretParams{
				Pass: ValidSecretPass_1,
			},
		)

		if serr != nil {
			t.Fatalf("Expected nil error, got %v", serr)
		}
	})

	t.Run("Get the updated secret fields", func(t *testing.T) {

		var secret *store.UserSecret
		var serr *store.StoreError

		secret, serr = st.GetUserSecret(ctx, ValidUserID_1)

		if serr != nil {
			t.Fatalf("Expected nil error, got %v", serr)
		}

		if secret.ID != ValidUserID_1 {
			t.Fatalf("Expected secret.ID '%s', got '%s'", ValidUserID_1, secret.ID)
		}

		if secret.Pass != ValidSecretPass_1 {
			t.Fatalf("Expected secret.Pass '%s', got '%s'", ValidSecretPass_1, secret.Pass)
		}
	})

	t.Run("Delete the secret", func(t *testing.T) {

		var serr *store.StoreError = st.DeleteUserSecret(ctx, ValidUserID_1)

		if serr != nil {
			t.Fatalf("Expected nil error, got %v", serr)
		}
	})

	t.Run("Delete the secret that doesn't exist", func(t *testing.T) {

		var serr *store.StoreError = st.DeleteUserSecret(ctx, ValidUserID_1)

		if serr == nil || serr.Code != store.ErrorCode_UserSecretNotFound {
			t.Fatalf("Expected store.ErrorCode_SecretNotFound error, got %v", serr)
		}
	})
}
