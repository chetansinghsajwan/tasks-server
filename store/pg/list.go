package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"tasks/store"

	"github.com/jackc/pgx/v5/pgconn"
)

func (st PostgresStore) ParseListID(id string) (store.ListID, error) {

	if len(strings.TrimSpace(id)) == 0 {
		return "", errors.New("list id must not be empty")
	}

	return store.ListID(id), nil
}

func (st PostgresStore) CreateList(ctx context.Context,
	args store.CreateListParams) *store.StoreError {

	const query = `
		insert into lists(id, owner_id)
		values ($1, $2)
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, query, args.ID, args.OwnerID); err != nil {

		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {

			store.PrintPgError(pgerr)
		}

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

func (st PostgresStore) GetList(ctx context.Context,
	id store.ListID, owner_id store.UserID) (*store.List, *store.StoreError) {

	const query = `
		select id, owner_id
		from lists
		where id = $1 and owner_id = $2
	`

	var list store.List
	var err error = st.Pool.QueryRow(ctx, query, id, owner_id).Scan(
		&list.ID, &list.OwnerID,
	)

	if err != nil {

		if errors.Unwrap(err) == sql.ErrNoRows {

			return nil, &store.StoreError{
				Code:         store.ErrListNotFoundCode,
				Msg:          fmt.Sprintf("list with id '%s', owner id '%s' not found", id, owner_id),
				WrappedError: err,
			}
		}

		return nil, &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return &list, nil
}

func (st PostgresStore) UpdateList(ctx context.Context, id store.ListID,
	owner_id store.UserID, args store.UpdateListParams) *store.StoreError {

	var queryBuilder strings.Builder
	queryBuilder.WriteString("update lists set")

	// First 2 positions is for id and owner_id
	var queryArgs = []any{id, owner_id}

	if args.ID.IsSome() {

		if len(queryArgs) > 3 {
			queryBuilder.WriteString(" ,")
		}

		queryBuilder.WriteString(fmt.Sprintf(" id = $%d", len(queryArgs)))
		queryArgs = append(queryArgs, args.ID.MustGet())
	}

	if args.OwnerID.IsSome() {

		if len(queryArgs) > 3 {
			queryBuilder.WriteString(" ,")
		}

		queryBuilder.WriteString(fmt.Sprintf(" owner_id = $%d", len(queryArgs)))
		queryArgs = append(queryArgs, args.OwnerID.MustGet())
	}

	// There are no updates, first 2 positions is for id and and owner_id
	if len(queryArgs) == 2 {
		return nil
	}

	queryBuilder.WriteString(" where id = $1 and owner_id = $2")

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, queryBuilder.String(), queryArgs...); err != nil {

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Update() {

		return &store.StoreError{
			Code:         store.ErrListNotFoundCode,
			Msg:          fmt.Sprintf("list with id '%s', owner id '%s' not found", id, owner_id),
			WrappedError: err,
		}
	}

	return nil
}

func (st PostgresStore) DeleteList(ctx context.Context,
	id store.ListID, owner_id store.UserID) *store.StoreError {

	const query = `
		delete from lists
		where id = $1 and owner_id = $2
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, query, id, owner_id); err != nil {

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Delete() {

		return &store.StoreError{
			Code:         store.ErrListNotFoundCode,
			Msg:          fmt.Sprintf("list with id '%s', owner id '%s' not found", id, owner_id),
			WrappedError: err,
		}
	}

	return nil
}
