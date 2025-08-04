package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"tasks/store"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresListAccessStore struct {
	Pool *pgxpool.Pool
}

func (las *PostgresListAccessStore) AddAccess(ctx context.Context, args store.ListAccess) *store.StoreError {

	const query = `
		insert into list_access (user_id, list_id, access)
		values ($1, $2, $3)
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = las.Pool.Exec(ctx, query, args.UserID, args.ListID, args.Access); err != nil {

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Insert() {

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return nil
}

func (las *PostgresListAccessStore) HasAccess(ctx context.Context, args store.ListAccess) (bool, *store.StoreError) {

	const query = `
		select access
		from list_access
		where user_id = $1 and list_id = $2 and access = $3
	`

	var access string
	var row = las.Pool.QueryRow(ctx, query, args.UserID, args.ListID, args.Access)
	var err = row.Scan(&access)

	if err != nil {

		if errors.Unwrap(err) == sql.ErrNoRows {
			return false, nil
		}

		return false, &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return true, nil
}

func (las *PostgresListAccessStore) RemoveAccesses(ctx context.Context, args store.RemoveListAccessParams) *store.StoreError {

	const query = `
		delete from list_access
		where user_id = $1 and list_id = $2 and access = $3
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = las.Pool.Exec(ctx, query, args.UserID, args.ListID, args.Access); err != nil {

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Delete() {

		return &store.StoreError{
			Code: store.ErrListAccessNotFound,
			Msg: fmt.Sprintf("list access '%s' for user '%s' and list '%s' not found",
				args.Access, args.UserID, args.ListID),
			WrappedError: err,
		}
	}

	return nil
}
