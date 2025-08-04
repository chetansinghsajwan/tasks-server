package store

import (
	"context"
	"fmt"
	"tasks/option"
	"tasks/utils"
	"time"
)

type TaskID int64

func NullTaskID() TaskID {
	return 0
}

func ParseTaskID(id string) (TaskID, error) {

	var parsedId int64
	var err error
	if parsedId, err = utils.ParseInt64(id); err != nil {
		return 0, err
	}

	return TaskID(parsedId), nil
}

func (id TaskID) String() string {
	return fmt.Sprintf("%d", id)
}

func (id TaskID) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%d", id)
}

type Task struct {
	ID          TaskID                   `db:"id" json:"id"`
	ListID      ListID                   `db:"list_id" json:"list_id"`
	Title       string                   `db:"title" json:"title"`
	Description option.Option[string]    `db:"description" json:"description"`
	Priority    option.Option[uint32]    `db:"priority" json:"priority"`
	DueDate     option.Option[time.Time] `db:"due_date" json:"due_date"`
	Assignee    option.Option[UserID]    `db:"assignee" json:"assignee"`
	Labels      []string                 `db:"labels" json:"labels"`
}

type CreateTaskParams struct {
	ListID      ListID                   `db:"list_id" json:"list_id"`
	Title       string                   `db:"title" json:"title"`
	Description option.Option[string]    `db:"description" json:"description"`
	Priority    option.Option[uint32]    `db:"priority" json:"priority"`
	DueDate     option.Option[time.Time] `db:"due_date" json:"due_date"`
	Assignee    option.Option[UserID]    `db:"assignee" json:"assignee"`
	Labels      []string                 `db:"labels" json:"labels"`
}

type UpdateTaskParams struct {
	ListID      option.Option[*ListID]    `db:"list_id" json:"list_id"`
	Title       option.Option[*string]    `db:"title" json:"title"`
	Description option.Option[*string]    `db:"description" json:"description"`
	Priority    option.Option[*uint32]    `db:"priority" json:"priority"`
	DueDate     option.Option[*time.Time] `db:"due_date" json:"due_date"`
	Assignee    option.Option[*UserID]    `db:"assignee" json:"assignee"`
	Labels      option.Option[[]string]   `db:"labels" json:"labels"`
}

type TaskStore interface {
	NullTaskID() TaskID
	ParseTaskID(id string) (TaskID, error)
	GetTask(ctx context.Context, id TaskID) (*Task, *StoreError)
	GetTasksWhere(ctx context.Context, where string, count uint, from uint) ([]Task, *StoreError)
	CreateTask(ctx context.Context, args CreateTaskParams) (TaskID, *StoreError)
	UpdateTask(ctx context.Context, id TaskID, args UpdateTaskParams) *StoreError
	DeleteTask(ctx context.Context, id TaskID) *StoreError
}
