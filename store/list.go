package store

import (
	"tasks/option"
)

type ListID string

type List struct {
	ID      ListID `db:"id" json:"id"`
	OwnerID UserID `db:"owner_id" json:"owner_id"`
}

type CreateListParams struct {
	ID      ListID `db:"id" json:"id"`
	OwnerID UserID `db:"owner_id" json:"owner_id"`
}

type UpdateListParams struct {
	ID      option.Option[*ListID] `db:"id" json:"id"`
	OwnerID option.Option[*UserID] `db:"owner_id" json:"owner_id"`
}
