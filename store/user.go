package store

import (
	"tasks/option"
)

type User struct {
	ID          string                `json:"id"`
	FullName    string                `json:"full_name"`
	DisplayName option.Option[string] `json:"display_name"`
	Email       string                `json:"email"`
}

type CreateUserParams struct {
	ID          string
	FullName    string
	DisplayName option.Option[string]
	Email       string
}

type UpdateUserParams struct {
	ID          option.Option[*string]
	Email       option.Option[*string]
	FullName    option.Option[*string]
	DisplayName option.Option[*string]
}
