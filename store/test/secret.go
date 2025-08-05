package test

import (
	"context"
	"tasks/option"
	"tasks/store"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestSecretStore(t *testing.T, st store.Store) {

	var ValidSecretID_0 = "myid_0"
	var ValidSecretValue_0 = "mypassword_0"
	var ValidSecretScope_0 = "user-login"

	var ValidSecretID_1 = "myid_1"
	var ValidSecretValue_1 = "mypassword_1"
	var ValidSecretScope_1 = "user-login"

	var ctx = context.Background()

	t.Run("Create secret with invalid id", func(t *testing.T) {

		var serr *store.StoreError = st.CreateSecret(ctx, store.CreateSecretParams{
			ID:    "  ",
			Scope: ValidSecretScope_0,
			Value: ValidSecretValue_0,
		})

		if serr == nil || serr.Code != store.ErrorCode_InvalidSecretIDFormat {
			t.Fatalf("Expected store.ErrorCode_InvalidSecretIDFormat error, got %v", serr)
		}
	})

	t.Run("Create secret with invalid scope", func(t *testing.T) {

		var serr *store.StoreError = st.CreateSecret(ctx, store.CreateSecretParams{
			ID:    ValidSecretID_0,
			Scope: "myscope",
			Value: ValidSecretValue_0,
		})

		if serr == nil || serr.Code != store.ErrorCode_InvalidSecretScope {
			t.Fatalf("Expected store.ErrorCode_InvalidSecretScope error, got %v", serr)
		}
	})

	t.Run("Create secret with invalid pass", func(t *testing.T) {

		var serr *store.StoreError = st.CreateSecret(ctx, store.CreateSecretParams{
			ID:    ValidSecretID_0,
			Scope: ValidSecretScope_0,
			Value: "  ",
		})

		if serr == nil || serr.Code != store.ErrorCode_InvalidSecretValueFormat {
			t.Fatalf("Expected store.ErrorCode_InvalidSecretValueFormat error, got %v", serr)
		}
	})

	t.Run("Create secret with valid fields", func(t *testing.T) {

		var serr *store.StoreError = st.CreateSecret(ctx, store.CreateSecretParams{
			ID:    ValidSecretID_0,
			Scope: ValidSecretScope_0,
			Value: ValidSecretValue_0,
		})

		if serr != nil {
			t.Fatalf("Expected nil error, got %v", serr)
		}
	})

	t.Run("Get secret that doesn't exist", func(t *testing.T) {

		var serr *store.StoreError

		_, serr = st.GetSecret(ctx, store.SecretKey{
			ID:    ValidSecretID_1,
			Scope: ValidSecretScope_0,
		})

		if serr == nil || serr.Code != store.ErrorCode_SecretNotFound {
			t.Fatalf("Expected store.ErrorCode_SecretNotFound error, got %v", serr)
		}
	})

	t.Run("Get valid secret", func(t *testing.T) {

		var secret *store.Secret
		var serr *store.StoreError

		secret, serr = st.GetSecret(ctx, store.SecretKey{
			ID:    ValidSecretID_0,
			Scope: ValidSecretScope_0,
		})

		if serr != nil {
			t.Fatalf("Expected nil error, got %v", serr)
		}

		if secret.ID != ValidSecretID_0 {
			t.Fatalf("Expected secret.ID '%s', got '%s'", ValidSecretID_0, secret.ID)
		}

		if secret.Scope != ValidSecretScope_0 {
			t.Fatalf("Expected secret.Scope '%s', got '%s'", ValidSecretScope_0, secret.Scope)
		}

		if secret.Value != ValidSecretValue_0 {
			t.Fatalf("Expected secret.Value '%s', got '%s'", ValidSecretValue_0, secret.Value)
		}
	})

	t.Run("Update secret with invalid id", func(t *testing.T) {

		var serr *store.StoreError = st.UpdateSecret(ctx,
			store.SecretKey{
				ID:    ValidSecretID_0,
				Scope: ValidSecretScope_0,
			},
			store.UpdateSecretParams{
				ID:    option.Some("  "),
				Scope: option.Some(ValidSecretScope_1),
				Value: option.Some(ValidSecretValue_1),
			},
		)

		if serr == nil || serr.Code != store.ErrorCode_InvalidSecretIDFormat {
			t.Fatalf("Expected store.ErrorCode_InvalidSecretIDFormat error, got %v", serr)
		}
	})

	t.Run("Update secret with invalid scope", func(t *testing.T) {

		var serr *store.StoreError = st.UpdateSecret(ctx,
			store.SecretKey{
				ID:    ValidSecretID_0,
				Scope: ValidSecretScope_0,
			},
			store.UpdateSecretParams{
				ID:    option.Some(ValidSecretID_1),
				Scope: option.Some("myscope"),
				Value: option.Some(ValidSecretValue_1),
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
				ID:    ValidSecretID_0,
				Scope: ValidSecretScope_0,
			},
			store.UpdateSecretParams{
				ID:    option.Some(ValidSecretID_1),
				Scope: option.Some(ValidSecretScope_1),
				Value: option.Some("  "),
			},
		)

		if serr == nil || serr.Code != store.ErrorCode_InvalidSecretValueFormat {
			store.PrintPgError(serr.WrappedError.(*pgconn.PgError))
			t.Fatalf("Expected store.ErrorCode_InvalidSecretValueFormat error, got %v", serr)
		}
	})

	t.Run("Update secret with valid fields", func(t *testing.T) {

		var serr *store.StoreError = st.UpdateSecret(ctx,
			store.SecretKey{
				ID:    ValidSecretID_0,
				Scope: ValidSecretScope_0,
			},
			store.UpdateSecretParams{
				ID:    option.Some(ValidSecretID_1),
				Scope: option.Some(ValidSecretScope_1),
				Value: option.Some(ValidSecretValue_1),
			},
		)

		if serr != nil {
			store.PrintPgError(serr.WrappedError.(*pgconn.PgError))
			t.Fatalf("Expected nil error, got %v", serr)
		}
	})

	t.Run("Get the updated secret fields", func(t *testing.T) {

		var secret *store.Secret
		var serr *store.StoreError

		secret, serr = st.GetSecret(ctx, store.SecretKey{
			ID:    ValidSecretID_1,
			Scope: ValidSecretScope_1,
		})

		if serr != nil {
			t.Fatalf("Expected nil error, got %v", serr)
		}

		if secret.ID != ValidSecretID_1 {
			t.Fatalf("Expected secret.ID '%s', got '%s'", ValidSecretID_1, secret.ID)
		}

		if secret.Scope != ValidSecretScope_1 {
			t.Fatalf("Expected secret.Scope '%s', got '%s'", ValidSecretScope_1, secret.Scope)
		}

		if secret.Value != ValidSecretValue_1 {
			t.Fatalf("Expected secret.Value '%s', got '%s'", ValidSecretValue_1, secret.Value)
		}
	})

	t.Run("Delete the secret", func(t *testing.T) {

		var serr *store.StoreError = st.DeleteSecret(ctx, store.SecretKey{
			ID:    ValidSecretID_1,
			Scope: ValidSecretScope_1,
		})

		if serr != nil {
			t.Fatalf("Expected nil error, got %v", serr)
		}
	})

	t.Run("Delete the secret that doesn't exist", func(t *testing.T) {

		var serr *store.StoreError = st.DeleteSecret(ctx, store.SecretKey{
			ID:    ValidSecretID_1,
			Scope: ValidSecretScope_1,
		})

		if serr == nil || serr.Code != store.ErrorCode_SecretNotFound {
			t.Fatalf("Expected store.ErrorCode_SecretNotFound error, got %v", serr)
		}
	})
}
