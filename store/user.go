package store

import (
	"errors"
	"strings"
	"tasks/option"
)

type UserID string

type User struct {
	ID          UserID                `json:"id"`
	FullName    string                `json:"full_name"`
	DisplayName option.Option[string] `json:"display_name"`
	Email       string                `json:"email"`
}

type CreateUserParams struct {
	ID          UserID
	FullName    string
	DisplayName option.Option[string]
	Email       string
}

type UpdateUserParams struct {
	ID          option.Option[*UserID]
	Email       option.Option[*string]
	FullName    option.Option[*string]
	DisplayName option.Option[*string]
}

func NullUserID() UserID {
	return ""
}

func ParseUserID(id string) (UserID, error) {

	if len(strings.TrimSpace(id)) == 0 {
		return "", errors.New("user id must not be empty")
	}

	return UserID(id), nil
}

func (id UserID) String() string {
	return string(id)
}
