package store

import (
	"tasks/option"
)

type ListAccessType string

const (
	ListAccessType_Owner        ListAccessType = "owner"
	ListAccessType_AddTask      ListAccessType = "add-task"
	ListAccessType_ReadTask     ListAccessType = "read-task"
	ListAccessType_WriteTask    ListAccessType = "write-task"
	ListAccessType_RemoveTask   ListAccessType = "remove-task"
	ListAccessType_AddAccess    ListAccessType = "add-access"
	ListAccessType_ReadAccess   ListAccessType = "read-access"
	ListAccessType_RemoveAccess ListAccessType = "remove-access"
)

type ListAccess struct {
	UserID string
	ListID uint64
	Access []ListAccessType
}

type GetListAccessParams struct {
	UserID string
	ListID uint64
}

type AddListAccessParams struct {
	UserID string
	ListID uint64
	Access []ListAccessType
}

type RemoveListAccessParams struct {
	UserID option.Option[string]
	ListID option.Option[uint64]
	Access option.Option[[]ListAccessType]
}
