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
	UserID UserID
	ListID ListID
	Access []ListAccessType
}

type GetListAccessParams struct {
	UserID UserID
	ListID ListID
}

type AddListAccessParams struct {
	UserID UserID
	ListID ListID
	Access []ListAccessType
}

type RemoveListAccessParams struct {
	UserID option.Option[UserID]
	ListID option.Option[ListID]
	Access option.Option[[]ListAccessType]
}
