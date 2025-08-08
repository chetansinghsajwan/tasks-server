package store

import (
	"context"
)

type Store interface {
	CreateUser(ctx context.Context, args CreateUserParams) *StoreError
	GetUser(ctx context.Context, id string) (*User, *StoreError)
	GetUsersWhere(ctx context.Context, where string, count uint, from uint) ([]User, *StoreError)
	UpdateUser(ctx context.Context, id string, args UpdateUserParams) *StoreError
	DeleteUser(ctx context.Context, id string) *StoreError

	CreateUserSecret(ctx context.Context, args CreateUserSecretParams) *StoreError
	GetUserSecret(ctx context.Context, id string) (*UserSecret, *StoreError)
	UpdateUserSecret(ctx context.Context, id string, args UpdateUserSecretParams) *StoreError
	DeleteUserSecret(ctx context.Context, id string) *StoreError

	CreateList(ctx context.Context, args CreateListParams) (uint64, *StoreError)
	GetList(ctx context.Context, id uint64) (*List, *StoreError)
	UpdateList(ctx context.Context, id uint64, args UpdateListParams) *StoreError
	DeleteList(ctx context.Context, id uint64) *StoreError

	GetTask(ctx context.Context, id TaskID) (*Task, *StoreError)
	CreateTask(ctx context.Context, args CreateTaskParams) (TaskID, *StoreError)
	UpdateTask(ctx context.Context, id TaskID, args UpdateTaskParams) *StoreError
	DeleteTask(ctx context.Context, id TaskID) *StoreError

	GetListAccess(ctx context.Context, args GetListAccessParams) (*ListAccess, *StoreError)
	AddListAccess(ctx context.Context, args AddListAccessParams) *StoreError
	RemoveListAccess(ctx context.Context, args RemoveListAccessParams) *StoreError
}
