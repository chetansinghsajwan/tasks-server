package store

import (
	"fmt"
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
	ID          TaskID     `db:"id" json:"id"`
	ListID      uint64     `db:"list_id" json:"list_id"`
	Title       string     `db:"title" json:"title"`
	Description *string    `db:"description" json:"description"`
	Priority    *uint32    `db:"priority" json:"priority"`
	DueDate     *time.Time `db:"due_date" json:"due_date"`
	Assignee    *string    `db:"assignee" json:"assignee"`
	Labels      []string   `db:"labels" json:"labels"`
}

type CreateTaskParams struct {
	ListID      uint64     `db:"list_id" json:"list_id"`
	Title       string     `db:"title" json:"title"`
	Description *string    `db:"description" json:"description"`
	Priority    *uint32    `db:"priority" json:"priority"`
	DueDate     *time.Time `db:"due_date" json:"due_date"`
	Assignee    *string    `db:"assignee" json:"assignee"`
	Labels      []string   `db:"labels" json:"labels"`
}

type UpdateTaskParams struct {
	ListID      **uint64    `db:"list_id" json:"list_id"`
	Title       **string    `db:"title" json:"title"`
	Description **string    `db:"description" json:"description"`
	Priority    **uint32    `db:"priority" json:"priority"`
	DueDate     **time.Time `db:"due_date" json:"due_date"`
	Assignee    **string    `db:"assignee" json:"assignee"`
	Labels      *[]string   `db:"labels" json:"labels"`
}
