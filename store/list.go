package store

import (
	"errors"
	"strings"
	"tasks/option"
)

type ListID string

type List struct {
	ID   ListID `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type CreateListParams struct {
	Name string `db:"name" json:"name"`
}

type UpdateListParams struct {
	Name option.Option[string] `db:"name" json:"name"`
}

func NullListID() ListID {
	return ""
}

func ParseListID(id string) (ListID, error) {

	if len(strings.TrimSpace(id)) == 0 {
		return "", errors.New("list id must not be empty")
	}

	return ListID(id), nil
}

func (id ListID) String() string {
	return string(id)
}
