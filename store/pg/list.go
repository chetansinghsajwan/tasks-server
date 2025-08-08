package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"tasks/errorcodes"
	"tasks/store"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (st PostgresStore) CreateList(ctx context.Context,
	args store.CreateListParams) (store.ListID, *store.StoreError) {

	const query = `
		insert into lists(name)
		values ($1)
		returning id
	`

	var row pgx.Row = st.Pool.QueryRow(ctx, query, args.Name)

	var listID store.ListID
	var err error
	if err = row.Scan(&listID); err != nil {

		return store.NullListID(), &store.StoreError{
			Code:         errorcodes.Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return listID, nil
}

func (st PostgresStore) GetList(ctx context.Context,
	id store.ListID) (*store.List, *store.StoreError) {

	const query = `
		select id, name
		from lists
		where id = $1
	`

	var list store.List
	var err error = st.Pool.QueryRow(ctx, query, id).Scan(
		&list.ID, &list.Name,
	)

	if err != nil {

		if errors.Unwrap(err) == sql.ErrNoRows {

			return nil, &store.StoreError{
				Code:         errorcodes.ListNotFound,
				Msg:          fmt.Sprintf("list with id '%s' not found", id),
				WrappedError: err,
			}
		}

		return nil, &store.StoreError{
			Code:         errorcodes.Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return &list, nil
}

func (st PostgresStore) UpdateList(ctx context.Context, id store.ListID, args store.UpdateListParams) *store.StoreError {

	const query = `
		update lists set
		name = $2
		where id = $1
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, query, id, args.Name); err != nil {

		return &store.StoreError{
			Code:         errorcodes.Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Update() {

		return &store.StoreError{
			Code:         errorcodes.ListNotFound,
			Msg:          fmt.Sprintf("list with id '%s' not found", id),
			WrappedError: err,
		}
	}

	return nil
}

func (st PostgresStore) DeleteList(ctx context.Context,
	id store.ListID) *store.StoreError {

	const query = `
		delete from lists
		where id = $1
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, query, id); err != nil {

		return &store.StoreError{
			Code:         errorcodes.Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Delete() {

		return &store.StoreError{
			Code:         errorcodes.ListNotFound,
			Msg:          fmt.Sprintf("list with id '%s' not found", id),
			WrappedError: err,
		}
	}

	return nil
}
