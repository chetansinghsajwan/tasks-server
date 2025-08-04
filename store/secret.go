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
