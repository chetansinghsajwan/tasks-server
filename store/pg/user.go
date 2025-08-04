package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"tasks/store"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const userIdFormatHint string = ""
const userEmailFormatHint string = ""
const userFullNameFormatHint string = ""
const userDisplayNameFormatHint string = ""

func (st PostgresStore) GetUser(ctx context.Context, id store.UserID) (*store.User, *store.StoreError) {

	const query = `
		SELECT id, email, full_name, display_name
		FROM users
		WHERE id = $1
	`

	var user store.User
	var err = st.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.FullName, &user.DisplayName,
	)

	if err != nil {

		if errors.Unwrap(err) == sql.ErrNoRows {

			return nil, &store.StoreError{
				Code:         store.ErrUserNotFoundCode,
				Msg:          fmt.Sprintf("user with id '%s' not found", id),
				WrappedError: err,
			}
		}

		return nil, &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return &user, nil
}

func (st PostgresStore) GetUsersWhere(ctx context.Context, where string, count uint, from uint) ([]store.User, *store.StoreError) {

	const query = `
		SELECT id, email, full_name, display_name
		FROM users
		WHERE $1
		OFFSET $2
		LIMIT $3
	`

	var err error
	var rows pgx.Rows
	rows, err = st.Pool.Query(ctx, query, where, from, count)

	if err != nil {

		return nil, &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	var users []store.User
	for rows.Next() {
		var user store.User

		if err = rows.Scan(&user.ID, &user.Email, &user.FullName, &user.DisplayName); err != nil {

			return nil, &store.StoreError{
				Code:         store.ErrUnknown,
				Msg:          "unknown error",
				WrappedError: err,
			}
		}

		users = append(users, user)
	}

	if err != nil {

		return nil, &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return users, nil
}

func (st PostgresStore) CreateUser(ctx context.Context, args store.CreateUserParams) *store.StoreError {

	const query = `
		INSERT INTO users (id, email, full_name, display_name)
		VALUES ($1, $2, $3, $4);
	`
	var cmd pgconn.CommandTag
	var err error
	cmd, err = st.Pool.Exec(ctx, query, args.ID, args.Email, args.FullName, args.DisplayName)

	if err != nil {

		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {

			if pgerr.Code == pgerrcode.NotNullViolation && pgerr.SchemaName == "public" &&
				pgerr.TableName == "users" && pgerr.ColumnName == "id" {

				return &store.StoreError{
					Code:         store.ErrUserIDNullCode,
					Msg:          "user id cannot be null",
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.CheckViolation && pgerr.ConstraintName == "users_id_validation" {

				return &store.StoreError{
					Code:         store.ErrUserIDFormatCode,
					Msg:          fmt.Sprintf("user id '%s' format is not correct. hint: %s", args.ID, userIdFormatHint),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.UniqueViolation && pgerr.ConstraintName == "users_pkey" {

				return &store.StoreError{
					Code:         store.ErrUserIDAlreadyExistsCode,
					Msg:          fmt.Sprintf("user id '%s' already exists", args.ID),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.NotNullViolation && pgerr.SchemaName == "public" &&
				pgerr.TableName == "users" && pgerr.ColumnName == "email" {

				return &store.StoreError{
					Code:         store.ErrUserEmailNullCode,
					Msg:          "user email cannot be null",
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.CheckViolation && pgerr.ConstraintName == "users_email_validation" {

				return &store.StoreError{
					Code:         store.ErrUserEmailFormatCode,
					Msg:          fmt.Sprintf("user email '%s' format is not correct. hint: %s", args.Email, userEmailFormatHint),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.UniqueViolation && pgerr.ConstraintName == "users_email_key" {

				return &store.StoreError{
					Code:         store.ErrUserEmailAlreadyExistsCode,
					Msg:          fmt.Sprintf("user email '%s' already exis", args.Email),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.CheckViolation && pgerr.ConstraintName == "users_full_name_validation" {

				return &store.StoreError{
					Code:         store.ErrUserFullNameFormatCode,
					Msg:          fmt.Sprintf("user full name '%s' format is not correct. hint: %s", args.FullName, userFullNameFormatHint),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.CheckViolation && pgerr.ConstraintName == "users_display_name_validation" {

				return &store.StoreError{
					Code:         store.ErrUserDisplayNameFormatCode,
					Msg:          fmt.Sprintf("user display name '%s' format is not correct. hint: %s", args.FullName, userDisplayNameFormatHint),
					WrappedError: err,
				}
			}
		}

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Insert() {

		return &store.StoreError{
			Code: store.ErrUnknown,
			Msg:  "user insertion failed",
		}
	}

	return nil
}

func (st PostgresStore) UpdateUser(ctx context.Context, id store.UserID, args store.UpdateUserParams) *store.StoreError {

	// Build the query
	var queryBuilder strings.Builder
	queryBuilder.WriteString("UPDATE users SET ")

	// $1 is set for 'id'
	var argIndex = 2
	var queryArgs = []any{id}

	if args.Email.IsSome() {
		queryBuilder.WriteString(fmt.Sprintf("email = %d", argIndex))
		queryArgs = append(queryArgs, args.Email)
		argIndex += 1
	}

	if args.FullName.IsSome() {
		queryBuilder.WriteString(fmt.Sprintf("full_name = %d", argIndex))
		queryArgs = append(queryArgs, args.FullName)
		argIndex += 1
	}

	if args.DisplayName.IsSome() {
		queryBuilder.WriteString(fmt.Sprintf("display_name = %d", argIndex))
		queryArgs = append(queryArgs, args.DisplayName)
		argIndex += 1
	}

	// There were no updates
	if argIndex == 2 {
		return nil
	}

	queryBuilder.WriteString(" WHERE id = $1")

	// Execute the query
	var cmd pgconn.CommandTag
	var err error
	cmd, err = st.Pool.Exec(ctx, queryBuilder.String(), queryArgs...)

	// Error handling
	if err != nil {

		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {

			if pgerr.Code == pgerrcode.CheckViolation && pgerr.ConstraintName == "users_id_validation" {

				return &store.StoreError{
					Code:         store.ErrUserIDFormatCode,
					Msg:          fmt.Sprintf("user id '%s' format is not correct. hint: %s", args.ID, userIdFormatHint),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.UniqueViolation && pgerr.ConstraintName == "users_pkey" {

				return &store.StoreError{
					Code:         store.ErrUserIDAlreadyExistsCode,
					Msg:          fmt.Sprintf("user id '%s' already exists", args.ID),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.NotNullViolation && pgerr.SchemaName == "public" &&
				pgerr.TableName == "users" && pgerr.ColumnName == "email" {

				return &store.StoreError{
					Code:         store.ErrUserEmailNullCode,
					Msg:          "user email cannot be null",
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.CheckViolation && pgerr.ConstraintName == "users_email_validation" {

				return &store.StoreError{
					Code:         store.ErrUserEmailFormatCode,
					Msg:          fmt.Sprintf("user email '%s' format is not correct. hint: %s", args.Email, userEmailFormatHint),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.UniqueViolation && pgerr.ConstraintName == "users_email_key" {

				return &store.StoreError{
					Code:         store.ErrUserEmailAlreadyExistsCode,
					Msg:          fmt.Sprintf("user email '%s' already exis", args.Email),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.CheckViolation && pgerr.ConstraintName == "users_full_name_validation" {

				return &store.StoreError{
					Code:         store.ErrUserFullNameFormatCode,
					Msg:          fmt.Sprintf("user full name '%s' format is not correct. hint: %s", args.FullName, userFullNameFormatHint),
					WrappedError: err,
				}
			}

			if pgerr.Code == pgerrcode.CheckViolation && pgerr.ConstraintName == "users_display_name_validation" {

				return &store.StoreError{
					Code:         store.ErrUserDisplayNameFormatCode,
					Msg:          fmt.Sprintf("user display name '%s' format is not correct. hint: %s", args.FullName, userDisplayNameFormatHint),
					WrappedError: err,
				}
			}
		}

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Update() {

		return &store.StoreError{
			Code: store.ErrUserNotFoundCode,
			Msg:  fmt.Sprintf("user '%s' not found", id),
		}
	}

	return nil
}

func (st PostgresStore) DeleteUser(ctx context.Context, id store.UserID) *store.StoreError {

	const query = `
		DELETE FROM users
		WHERE id = $0
	`

	var cmd pgconn.CommandTag
	var err error
	cmd, err = st.Pool.Exec(ctx, query, id)

	if err != nil {

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Delete() {

		return &store.StoreError{
			Code: store.ErrUserNotFoundCode,
			Msg:  fmt.Sprintf("user '%s' not found", id),
		}
	}

	return nil
}
