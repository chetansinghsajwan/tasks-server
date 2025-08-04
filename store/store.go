package store

import (
	"context"
)

type Store interface {
	CreateUser(ctx context.Context, args CreateUserParams) *StoreError
	GetUser(ctx context.Context, id UserID) (*User, *StoreError)
	GetUsersWhere(ctx context.Context, where string, count uint, from uint) ([]User, *StoreError)
	UpdateUser(ctx context.Context, id UserID, args UpdateUserParams) *StoreError
	DeleteUser(ctx context.Context, id UserID) *StoreError

	CreateList(ctx context.Context, args CreateListParams) (ListID, *StoreError)
	GetList(ctx context.Context, id ListID) (*List, *StoreError)
	UpdateList(ctx context.Context, id ListID, args UpdateListParams) *StoreError
	DeleteList(ctx context.Context, id ListID) *StoreError

	CreateSecret(ctx context.Context, args CreateSecretParams) *StoreError
	GetSecret(ctx context.Context, key SecretKey) (*Secret, *StoreError)
	UpdateSecret(ctx context.Context, key SecretKey, args UpdateSecretParams) *StoreError
	DeleteSecret(ctx context.Context, key SecretKey) *StoreError

	GetTask(ctx context.Context, id TaskID) (*Task, *StoreError)
	CreateTask(ctx context.Context, args CreateTaskParams) (TaskID, *StoreError)
	UpdateTask(ctx context.Context, id TaskID, args UpdateTaskParams) *StoreError
	DeleteTask(ctx context.Context, id TaskID) *StoreError

	GetListAccess(ctx context.Context, args GetListAccessParams) (*ListAccess, *StoreError)
	AddListAccess(ctx context.Context, args AddListAccessParams) *StoreError
	RemoveListAccess(ctx context.Context, args RemoveListAccessParams) *StoreError
}
