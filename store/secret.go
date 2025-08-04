package store

import (
	"context"
	"tasks/option"
)

type SecretKey struct {
	ID    string
	Scope string
}

type Secret struct {
	ID    string
	Scope string
	Pass  string
}

type CreateSecretParams struct {
	ID    string
	Scope string
	Pass  string
}

type UpdateSecretParams struct {
	ID    option.Option[string]
	Scope option.Option[string]
	Pass  option.Option[string]
}

type SecretStore interface {
	CreateSecret(ctx context.Context, args CreateSecretParams) *StoreError
	GetSecret(ctx context.Context, key SecretKey) (*Secret, *StoreError)
	UpdateSecret(ctx context.Context, key SecretKey, args UpdateSecretParams) *StoreError
	DeleteSecret(ctx context.Context, key SecretKey) *StoreError
}
