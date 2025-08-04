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

func (st PostgresStore) CreateSecret(ctx context.Context, args store.CreateSecretParams) *store.StoreError {

	const query = `
		insert into secrets (id, scope, pass)
		values ($1, $2, $3)
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, query, args.ID, args.Scope, args.Pass); err != nil {

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

func (st PostgresStore) GetSecret(ctx context.Context, key store.SecretKey) (*store.Secret, *store.StoreError) {

	const query = `
		select key, scope, pass
		from secrets
		where key = $1 and scope = $2
	`

	var secret store.Secret
	var row = st.Pool.QueryRow(ctx, query, key.ID, key.Scope)
	var err = row.Scan(&secret.ID, &secret.Scope, &secret.Pass)

	if err != nil {

		if errors.Unwrap(err) == sql.ErrNoRows {

			return nil, &store.StoreError{
				Code:         store.ErrSecretNotFound,
				Msg:          fmt.Sprintf("secret with id '%s', scope '%s' not found", key.ID, key.Scope),
				WrappedError: err,
			}
		}

		return nil, &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return &secret, nil
}

func (st PostgresStore) UpdateSecret(ctx context.Context, key store.SecretKey, args store.UpdateSecretParams) *store.StoreError {

	var queryBuilder strings.Builder
	queryBuilder.WriteString("update secrets set")

	// first 2 positions is for id and scope
	var queryArgs = []any{key.ID, key.Scope}

	if args.ID.IsSome() {

		if len(queryArgs) > 3 {
			queryBuilder.WriteString(" ,")
		}

		queryBuilder.WriteString(fmt.Sprintf(" id = $%d", len(queryArgs)))
		queryArgs = append(queryArgs, args.ID)
	}

	if args.Scope.IsSome() {

		if len(queryArgs) > 3 {
			queryBuilder.WriteString(" ,")
		}

		queryBuilder.WriteString(fmt.Sprintf(" scope = $%d", len(queryArgs)))
		queryArgs = append(queryArgs, args.Scope)
	}

	if args.Pass.IsSome() {

		if len(queryArgs) > 3 {
			queryBuilder.WriteString(" ,")
		}

		queryBuilder.WriteString(fmt.Sprintf(" pass = $%d", len(queryArgs)))
		queryArgs = append(queryArgs, args.Pass)
	}

	// There are no updates, the 2 args are id and scope
	if len(queryArgs) == 2 {
		return nil
	}

	queryBuilder.WriteString(" where id = $1 and scope = $2")

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
			Code:         store.ErrSecretNotFound,
			Msg:          fmt.Sprintf("secret with id '%s', scope '%s' not found", key.ID, key.Scope),
			WrappedError: err,
		}
	}

	return nil
}

func (st PostgresStore) DeleteSecret(ctx context.Context, key store.SecretKey) *store.StoreError {

	const query = `
		delete from secrets
		where key = $1 and scope = $2
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, query, key.ID, key.Scope); err != nil {

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Delete() {

		return &store.StoreError{
			Code:         store.ErrSecretNotFound,
			Msg:          fmt.Sprintf("secret with id '%s', key '%s' not found", key.ID, key.Scope),
			WrappedError: err,
		}
	}

	return nil
}
