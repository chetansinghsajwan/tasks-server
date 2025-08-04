package store

import (
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
