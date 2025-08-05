package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"tasks/store"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const invalidSecretIDFormatHint = ""
const invalidSecretValueFormatHint = ""
const invalidSecretScopeFormatHint = ""

func (st PostgresStore) CreateSecret(ctx context.Context, args store.CreateSecretParams) *store.StoreError {

	const query = `
		insert into secrets (id, scope, value)
		values ($1, $2, $3)
	`

	var err error
	if _, err = st.Pool.Exec(ctx, query, args.ID, args.Scope, args.Value); err != nil {

		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {

			if pgerr.Code == pgerrcode.CheckViolation &&
				pgerr.ConstraintName == "secrets_id_validation" {

				return &store.StoreError{
					Code:         store.ErrorCode_InvalidSecretIDFormat,
					Msg:          fmt.Sprintf("secret id '%s' format is not correct. hint: %s", args.ID, invalidSecretIDFormatHint),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.CheckViolation &&
				pgerr.ConstraintName == "secrets_value_validation" {

				return &store.StoreError{
					Code:         store.ErrorCode_InvalidSecretValueFormat,
					Msg:          fmt.Sprintf("secret value '%s' format is not correct. hint: %s", args.Value, invalidSecretValueFormatHint),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.InvalidTextRepresentation &&
				strings.HasPrefix(pgerr.Message, "invalid input value for enum secret_scopes:") {

				return &store.StoreError{
					Code:         store.ErrorCode_InvalidSecretScope,
					Msg:          fmt.Sprintf("secret scope '%s' value is invalid, hint: %s", args.Scope, invalidSecretScopeFormatHint),
					WrappedError: err,
				}
			}
		}

		return &store.StoreError{
			Code:         store.ErrorCode_Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return nil
}

func (st PostgresStore) GetSecret(ctx context.Context, key store.SecretKey) (*store.Secret, *store.StoreError) {

	const query = `
		select id, scope, value
		from secrets
		where id = $1 and scope = $2
	`

	var secret store.Secret
	var row = st.Pool.QueryRow(ctx, query, key.ID, key.Scope)
	var err = row.Scan(&secret.ID, &secret.Scope, &secret.Value)

	if err != nil {

		if errors.Unwrap(err) == sql.ErrNoRows {

			return nil, &store.StoreError{
				Code:         store.ErrorCode_SecretNotFound,
				Msg:          fmt.Sprintf("secret with id '%s', scope '%s' not found", key.ID, key.Scope),
				WrappedError: err,
			}
		}

		return nil, &store.StoreError{
			Code:         store.ErrorCode_Unknown,
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

		queryBuilder.WriteString(" id = $3")
		queryArgs = append(queryArgs, args.ID.MustGet())
	}

	if args.Scope.IsSome() {

		if len(queryArgs) > 2 {
			queryBuilder.WriteString(",")
		}

		queryBuilder.WriteString(fmt.Sprintf(" scope = $%d", len(queryArgs)+1))
		queryArgs = append(queryArgs, args.Scope.MustGet())
	}

	if args.Value.IsSome() {

		if len(queryArgs) > 2 {
			queryBuilder.WriteString(",")
		}

		queryBuilder.WriteString(fmt.Sprintf(" value = $%d", len(queryArgs)+1))
		queryArgs = append(queryArgs, args.Value.MustGet())
	}

	// There are no updates, the 2 args are id and scope
	if len(queryArgs) == 2 {
		return nil
	}

	queryBuilder.WriteString(" where id = $1 and scope = $2")

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, queryBuilder.String(), queryArgs...); err != nil {

		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {

			if pgerr.Code == pgerrcode.CheckViolation &&
				pgerr.ConstraintName == "secrets_id_validation" {

				return &store.StoreError{
					Code:         store.ErrorCode_InvalidSecretIDFormat,
					Msg:          fmt.Sprintf("secret id '%s' format is not correct. hint: %s", args.ID, invalidSecretIDFormatHint),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.CheckViolation &&
				pgerr.ConstraintName == "secrets_value_validation" {

				return &store.StoreError{
					Code:         store.ErrorCode_InvalidSecretValueFormat,
					Msg:          fmt.Sprintf("secret value '%s' format is not correct. hint: %s", args.Value, invalidSecretValueFormatHint),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.InvalidTextRepresentation &&
				strings.HasPrefix(pgerr.Message, "invalid input value for enum secret_scopes:") {

				return &store.StoreError{
					Code:         store.ErrorCode_InvalidSecretScope,
					Msg:          fmt.Sprintf("secret scope '%s' value is invalid, hint: %s", args.Scope, invalidSecretScopeFormatHint),
					WrappedError: err,
				}
			}
		}

		return &store.StoreError{
			Code:         store.ErrorCode_Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if cmd.RowsAffected() == 0 {

		return &store.StoreError{
			Code:         store.ErrorCode_SecretNotFound,
			Msg:          fmt.Sprintf("secret with id '%s', scope '%s' not found", key.ID, key.Scope),
			WrappedError: err,
		}
	}

	return nil
}

func (st PostgresStore) DeleteSecret(ctx context.Context, key store.SecretKey) *store.StoreError {

	const query = `
		delete from secrets
		where id = $1 and scope = $2
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = st.Pool.Exec(ctx, query, key.ID, key.Scope); err != nil {

		return &store.StoreError{
			Code:         store.ErrorCode_Unknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if cmd.RowsAffected() == 0 {

		return &store.StoreError{
			Code:         store.ErrorCode_SecretNotFound,
			Msg:          fmt.Sprintf("secret with id '%s', key '%s' not found", key.ID, key.Scope),
			WrappedError: err,
		}
	}

	return nil
}
