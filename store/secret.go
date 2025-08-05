package store

import (
	"tasks/option"
)

type SecretKey struct {
	ID    string
	Scope string
}

type Secret struct {
	ID    string
	Scope string
	Value string
}

type CreateSecretParams struct {
	ID    string
	Scope string
	Value string
}

type UpdateSecretParams struct {
	ID    option.Option[string]
	Scope option.Option[string]
	Value option.Option[string]
}
