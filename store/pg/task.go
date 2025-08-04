package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"tasks/store"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTaskStore struct {
	Pool *pgxpool.Pool
}

func (ts *PostgresTaskStore) CreateTask(ctx context.Context,
	args store.CreateTaskParams) (store.TaskID, *store.StoreError) {

	var query = `
		insert into tasks(list_id, title, description, priority, due_date, assignee, labels)
		values ($1, $2, $3, $4, $5, $6, $7)
		returning id
	`

	var taskId store.TaskID
	var row = ts.Pool.QueryRow(ctx, query, args.ListID, args.Title, args.Description, args.Priority, args.DueDate, args.Labels)
	var err = row.Scan(&taskId)

	if err != nil {

		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {

			store.PrintPgError(pgerr)
		}

		return store.NullTaskID(), &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return taskId, nil
}

func (ts *PostgresTaskStore) GetTask(ctx context.Context,
	id store.TaskID) (*store.Task, *store.StoreError) {

	var query = `
		select id, list_id, title, description, priority, due_date, assignee, labels
		from tasks
		where id = $1
	`

	var task store.Task
	var err error = ts.Pool.QueryRow(ctx, query, id).Scan(
		&task.ID, &task.ListID, &task.Title, &task.Description, &task.Priority, &task.DueDate, &task.Assignee, &task.Labels,
	)

	if err != nil {

		if errors.Unwrap(err) == sql.ErrNoRows {

			return nil, &store.StoreError{
				Code:         store.ErrTaskNotFoundCode,
				Msg:          fmt.Sprintf("task with id '%v' not found", id),
				WrappedError: err,
			}
		}

		return nil, &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	return &task, nil
}

func (ts *PostgresTaskStore) UpdateTask(ctx context.Context, id store.TaskID, args store.UpdateTaskParams) *store.StoreError {

	var queryBuilder strings.Builder
	queryBuilder.WriteString("update tasks set")

	// 1 is for id
	var queryArgIndex uint8 = 2
	var queryArgs = []any{id}

	if args.ListID.IsSome() {
		queryBuilder.WriteString(" list_id = $")
		queryBuilder.WriteString(string(queryArgIndex))
		queryArgs = append(queryArgs, args.ListID.MustGet())

		queryArgIndex += 1
	}

	if args.Title.IsSome() {
		queryBuilder.WriteString(" title = $")
		queryBuilder.WriteString(string(queryArgIndex))
		queryArgs = append(queryArgs, args.Title.MustGet())

		queryArgIndex += 1
	}

	if args.Description.IsSome() {
		queryBuilder.WriteString(" description = $")
		queryBuilder.WriteString(string(queryArgIndex))
		queryArgs = append(queryArgs, args.Description.MustGet())

		queryArgIndex += 1
	}

	if args.Priority.IsSome() {
		queryBuilder.WriteString(" priority = $")
		queryBuilder.WriteString(string(queryArgIndex))
		queryArgs = append(queryArgs, args.Priority.MustGet())

		queryArgIndex += 1
	}

	if args.DueDate.IsSome() {
		queryBuilder.WriteString(" due_date = $")
		queryBuilder.WriteString(string(queryArgIndex))
		queryArgs = append(queryArgs, args.DueDate.MustGet())

		queryArgIndex += 1
	}

	if args.Assignee.IsSome() {
		queryBuilder.WriteString(" assignee = $")
		queryBuilder.WriteString(string(queryArgIndex))
		queryArgs = append(queryArgs, args.Assignee.MustGet())

		queryArgIndex += 1
	}

	if args.Labels.IsSome() {
		queryBuilder.WriteString(" owner_id = $")
		queryBuilder.WriteString(string(queryArgIndex))
		queryArgs = append(queryArgs, args.Labels.MustGet())

		queryArgIndex += 1
	}

	// There are no updates
	if queryArgIndex == 2 {
		return nil
	}

	queryBuilder.WriteString(" where id = $1")

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = ts.Pool.Exec(ctx, queryBuilder.String(), queryArgs...); err != nil {

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Update() {

		return &store.StoreError{
			Code:         store.ErrTaskNotFoundCode,
			Msg:          fmt.Sprintf("task with id '%v' not found", id),
			WrappedError: err,
		}
	}

	return nil
}

func (ts *PostgresTaskStore) DeleteTask(ctx context.Context,
	id store.TaskID) *store.StoreError {

	var query = `
		delete from tasks
		where id = $1
	`

	var cmd pgconn.CommandTag
	var err error
	if cmd, err = ts.Pool.Exec(ctx, query, id); err != nil {

		return &store.StoreError{
			Code:         store.ErrUnknown,
			Msg:          "unknown error",
			WrappedError: err,
		}
	}

	if !cmd.Delete() {

		return &store.StoreError{
			Code:         store.ErrTaskNotFoundCode,
			Msg:          fmt.Sprintf("task with id '%v' not found", id),
			WrappedError: err,
		}
	}

	return nil
}
