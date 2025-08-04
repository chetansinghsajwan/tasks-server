package store

import (
	"context"
	"tasks/option"
)

type UserID string

type User struct {
	ID          string        `json:"id"`
	FullName    string        `json:"full_name"`
	DisplayName option.String `json:"display_name"`
	Email       string        `json:"email"`
}

type CreateUserParams struct {
	ID          string
	FullName    string
	DisplayName option.String
	Email       string
}

type UpdateUserParams struct {
	ID          option.StringPtr
	Email       option.StringPtr
	FullName    option.StringPtr
	DisplayName option.StringPtr
}

type UserStore interface {
	ParseUserID(id string) (UserID, error)
	CreateUser(ctx context.Context, args CreateUserParams) *StoreError
	GetUser(ctx context.Context, id UserID) (*User, *StoreError)
	GetUsersWhere(ctx context.Context, where string, count uint, from uint) ([]User, *StoreError)
	UpdateUser(ctx context.Context, id UserID, args UpdateUserParams) *StoreError
	DeleteUser(ctx context.Context, id UserID) *StoreError
}
