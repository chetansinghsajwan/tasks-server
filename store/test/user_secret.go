package test

import (
	"context"
	"tasks/option"
	"tasks/store"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestUserSecretStore(t *testing.T, st store.Store) {

	var ValidUserID_0, _ = store.ParseUserID("myid_0")
	var ValidSecretPass_0 = "mypassword_0"
	var ValidSecretScope_0 = "user-login"

	var ValidUserID_1, _ = store.ParseUserID("myid_1")
	var ValidSecretPass_1 = "mypassword_1"
	var ValidSecretScope_1 = "user-login"

	var ctx = context.Background()

	t.Run("Create secret with invalid id", func(t *testing.T) {

		var serr *store.StoreError = st.CreateUserSecret(ctx, store.CreateUserSecretParams{
			ID:   "  ",
			Pass: ValidSecretPass_0,
		})

		if serr == nil || serr.Code != store.ErrorCode_InvalidSecretIDFormat {
			t.Fatalf("Expected store.ErrorCode_InvalidSecretIDFormat error, got %v", serr)
		}
	})

	t.Run("Create secret with invalid scope", func(t *testing.T) {

		var serr *store.StoreError = st.CreateUserSecret(ctx, store.CreateUserSecretParams{
			ID:   ValidUserID_0,
			Pass: ValidSecretPass_0,
		})

		if serr == nil || serr.Code != store.ErrorCode_InvalidSecretScope {
			t.Fatalf("Expected store.ErrorCode_InvalidSecretScope error, got %v", serr)
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

		_, serr = st.GetSecret(ctx, ValidUserID_1)

		if serr == nil || serr.Code != store.ErrorCode_UserSecretNotFound {
			t.Fatalf("Expected store.ErrorCode_SecretNotFound error, got %v", serr)
		}
	})

	t.Run("Get valid secret", func(t *testing.T) {

		var secret *store.UserSecret
		var serr *store.StoreError

		secret, serr = st.GetSecret(ctx, ValidUserID_0)

		if serr != nil {
			t.Fatalf("Expected nil error, got %v", serr)
		}

		if secret.ID != ValidUserID_0 {
			t.Fatalf("Expected secret.ID '%s', got '%s'", ValidUserID_0, secret.ID)
		}

		if secret.Value != ValidSecretPass_0 {
			t.Fatalf("Expected secret.Value '%s', got '%s'", ValidSecretPass_0, secret.Value)
		}
	})

	t.Run("Update secret with invalid id", func(t *testing.T) {

		var serr *store.StoreError = st.UpdateSecret(ctx, ValidUserID_0,
			store.UpdateUserSecretParams{
				Pass: ValidSecretPass_1,
			},
		)

		if serr == nil || serr.Code != store.ErrorCode_InvalidSecretIDFormat {
			t.Fatalf("Expected store.ErrorCode_InvalidSecretIDFormat error, got %v", serr)
		}
	})

	t.Run("Update secret with invalid scope", func(t *testing.T) {

		var serr *store.StoreError = st.UpdateSecret(ctx,
			store.SecretKey{
				ID: ValidUserID_0,
			},
			store.UpdateUserSecretParams{
				ID:   option.Some(ValidUserID_1),
				Pass: option.Some(ValidSecretPass_1),
			},
		)

		if serr == nil || serr.Code != store.ErrorCode_InvalidSecretScope {
			store.PrintPgError(serr.WrappedError.(*pgconn.PgError))
			t.Fatalf("Expected store.ErrorCode_InvalidSecretScope error, got %v", serr)
		}
	})

	t.Run("Update secret with invalid pass", func(t *testing.T) {

		var serr *store.StoreError = st.UpdateSecret(ctx,
			store.SecretKey{
				ID: ValidUserID_0,
			},
			store.UpdateUserSecretParams{
				ID:   option.Some(ValidUserID_1),
				Pass: option.Some("  "),
			},
		)

		if serr == nil || serr.Code != store.ErrorCode_InvalidUserSecretPassFormat {
			store.PrintPgError(serr.WrappedError.(*pgconn.PgError))
			t.Fatalf("Expected store.ErrorCode_InvalidSecretValueFormat error, got %v", serr)
		}
	})

	t.Run("Update secret with valid fields", func(t *testing.T) {

		var serr *store.StoreError = st.UpdateSecret(ctx,
			store.SecretKey{
				ID: ValidUserID_0,
			},
			store.UpdateUserSecretParams{
				ID:   option.Some(ValidUserID_1),
				Pass: option.Some(ValidSecretPass_1),
			},
		)

		if serr != nil {
			store.PrintPgError(serr.WrappedError.(*pgconn.PgError))
			t.Fatalf("Expected nil error, got %v", serr)
		}
	})

	t.Run("Get the updated secret fields", func(t *testing.T) {

		var secret *store.UserSecret
		var serr *store.StoreError

		secret, serr = st.GetSecret(ctx, store.SecretKey{
			ID: ValidUserID_1,
		})

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
