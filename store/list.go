package store

import (
	"context"
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

type ListStore interface {
	ParseListID(id string) (ListID, error)
	CreateList(ctx context.Context, args CreateListParams) *StoreError
	GetList(ctx context.Context, id ListID, owner_id UserID) (*List, *StoreError)
	UpdateList(ctx context.Context, id ListID, owner_id UserID, args UpdateListParams) *StoreError
	DeleteList(ctx context.Context, id ListID, owner_id UserID) *StoreError
}
