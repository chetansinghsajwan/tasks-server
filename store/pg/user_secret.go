package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"tasks/errorcodes"
	"tasks/store"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const invalidUserSecretPassFormatHint string = ""

func (st PostgresStore) CreateUserSecret(ctx context.Context, args store.CreateUserSecretParams) *store.StoreError {

	const query = `
		insert into user_secrets (id, value)
		values ($1, $2)
	`

	var err error
	if _, err = st.Pool.Exec(ctx, query, args.ID, args.Pass); err != nil {

		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {

			if pgerr.Code == pgerrcode.CheckViolation &&
				pgerr.ConstraintName == "user_secrets_value_format_check" {

				return &store.StoreError{
					Code:         errorcodes.InvalidUserSecretPassFormat,
					Msg:          fmt.Sprintf("user secret value '%s' format is not correct. hint: %s", args.Pass, invalidUserSecretPassFormatHint),
					WrappedError: err,
				}
			}
		}

		return &store.StoreError{
			Code:         errorcodes.Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return nil
}

func (st PostgresStore) GetUserSecret(ctx context.Context, id string) (*store.UserSecret, *store.StoreError) {

	const query = `
		select id, value
		from user_secrets
		where id = $1
	`

	var secret store.UserSecret
	var row = st.Pool.QueryRow(ctx, query, id)
	var err = row.Scan(&secret.ID, &secret.Pass)

	if err != nil {

		if errors.Unwrap(err) == sql.ErrNoRows {

			return nil, &store.StoreError{
				Code:         errorcodes.UserSecretNotFound,
				Msg:          fmt.Sprintf("user secret with id '%s' not found", id),
				WrappedError: err,
			}
		}

		return nil, &store.StoreError{
			Code:         errorcodes.Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return &secret, nil
}

func (st PostgresStore) UpdateUserSecret(ctx context.Context, id string, args store.UpdateUserSecretParams) *store.StoreError {

	const query = `
		update users set
		value = $2
		where id = $1
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, query, id, args.Pass); err != nil {

		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {

			if pgerr.Code == pgerrcode.CheckViolation &&
				pgerr.ConstraintName == "user_secrets_value_format_check" {

				return &store.StoreError{
					Code:         errorcodes.InvalidUserSecretPassFormat,
					Msg:          fmt.Sprintf("user secret value '%s' format is not correct. hint: %s", args.Pass, invalidUserSecretPassFormatHint),
					WrappedError: err,
				}
			}
		}

		return &store.StoreError{
			Code:         errorcodes.Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if cmd.RowsAffected() == 0 {

		return &store.StoreError{
			Code:         errorcodes.UserSecretNotFound,
			Msg:          fmt.Sprintf("user secret with id '%s' not found", id),
			WrappedError: err,
		}
	}

	return nil
}

func (st PostgresStore) DeleteUserSecret(ctx context.Context, id string) *store.StoreError {

	const query = `
		delete from user_secrets
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

	if cmd.RowsAffected() == 0 {

		return &store.StoreError{
			Code:         errorcodes.UserSecretNotFound,
			Msg:          fmt.Sprintf("user secret with id '%s' not found", id),
			WrappedError: err,
		}
	}

	return nil
}
