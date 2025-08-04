package pg

import (
	"context"
	"fmt"
	"strings"
	"tasks/store"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (st PostgresStore) AddListAccess(ctx context.Context, args store.AddListAccessParams) *store.StoreError {

	// There are no accesses to add
	if len(args.Access) == 0 {
		return nil
	}

	// Build the query
	var queryArgs = []any{args.UserID, args.ListID}
	var queryBuilder strings.Builder

	queryBuilder.WriteString("insert into list_access (user_id, list_id, access) values ")

	// Add the value placeholders to the query
	for i, access := range args.Access {

		if i > 0 {
			queryBuilder.WriteString(", ")
		}

		queryBuilder.WriteString(fmt.Sprintf("($1, $2, $%s)", i+3))
		queryArgs = append(queryArgs, access)
	}

	// Add the accesses
	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, queryBuilder.String(), queryArgs...); err != nil {

		return &store.StoreError{
			Code:         store.ErrorCode_Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Insert() {

		return &store.StoreError{
			Code:         store.ErrorCode_Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return nil
}

func (st PostgresStore) GetListAccess(ctx context.Context, args store.GetListAccessParams) (*store.ListAccess, *store.StoreError) {

	const query = `
		select access
		from list_access
		where user_id = $1 and list_id = $2
	`

	var rows pgx.Rows
	var err error

	// Get the accesses
	if rows, err = st.Pool.Query(ctx, query, args.UserID, args.ListID); err != nil {

		return nil, &store.StoreError{
			Code:         store.ErrorCode_Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	defer rows.Close()

	// Read the accesses
	var accesses []store.ListAccessType
	for rows.Next() {

		var access store.ListAccessType
		if err = rows.Scan(&access); err != nil {

			return nil, &store.StoreError{
				Code:         store.ErrorCode_Unknown,
				Msg:          "unknown error",
				WrappedError: err,
			}
		}

		accesses = append(accesses, access)
	}

	// Check for errors in rows
	if err = rows.Err(); err != nil {

		return nil, &store.StoreError{
			Code:         store.ErrorCode_Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return &store.ListAccess{
		UserID: args.UserID,
		ListID: args.ListID,
		Access: accesses,
	}, nil
}

func (st PostgresStore) RemoveListAccess(ctx context.Context, args store.RemoveListAccessParams) *store.StoreError {

	// There is nothing to remove
	if args.UserID.IsNone() && args.ListID.IsNone() && args.Access.IsNone() {
		return nil
	}

	// Build the query
	var queryArgs []any
	var queryBuilder strings.Builder
	queryBuilder.WriteString("delete from list_access where ")

	if args.UserID.IsSome() {
		queryBuilder.WriteString("user_id = $1")
		queryArgs = append(queryArgs, args.UserID.MustGet())
	}

	if args.ListID.IsSome() {

		if len(queryArgs) != 0 {
			queryBuilder.WriteString(", ")
		}

		queryBuilder.WriteString(fmt.Sprintf("list_id = $%s", len(queryArgs)+1))
		queryArgs = append(queryArgs, args.ListID.MustGet())
	}

	if args.Access.IsSome() {

		if len(queryArgs) != 0 {
			queryBuilder.WriteString(", ")
		}

		queryBuilder.WriteString(fmt.Sprintf("access = any($%s::list_access_type[])", len(queryArgs)+1))
		queryArgs = append(queryArgs, args.Access.MustGet())
	}

	// Delete the accesses
	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, queryBuilder.String(), queryArgs...); err != nil {

		return &store.StoreError{
			Code:         store.ErrorCode_Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Delete() {

		return &store.StoreError{
			Code: store.ErrorCode_Unknown,
			Msg:  "unknown error",
		}
	}

	return nil
}
