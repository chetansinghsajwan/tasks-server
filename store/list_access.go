package store

import (
	"context"
	"tasks/option"
)

type ListAccessType string

const (
	ListAccessType_Owner  ListAccessType = "owner"
	ListAccessType_Read   ListAccessType = "read"
	ListAccessType_Write  ListAccessType = "write"
	ListAccessType_Update ListAccessType = "update"
	ListAccessType_Delete ListAccessType = "delete"
)

type ListAccess struct {
	UserID UserID
	ListID ListID
	Access ListAccessType
}

type RemoveListAccessParams struct {
	UserID option.Option[UserID]
	ListID option.Option[ListID]
	Access option.Option[ListAccessType]
}

type ListAccessStore interface {
	AddAccess(ctx context.Context, args ListAccess) *StoreError
	HasAccess(ctx context.Context, args ListAccess) (bool, *StoreError)
	RemoveAccesses(ctx context.Context, args RemoveListAccessParams) *StoreError
}
