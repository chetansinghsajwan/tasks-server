package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

type TaskID uint64

func ParseTaskID(s string) (TaskID, error) {
	parsed, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		return 0, err
	}

	return TaskID(parsed), nil
}

func (id TaskID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

type Task struct {
	ID          TaskID
	Title       string
	Description *string
	Priority    *uint8
	DueDate     *time.Time
	Assignee    *string
	Labels      []string
}

type TaskCreate struct {
	Title       string
	Description *string
	Priority    *uint8
	DueDate     *time.Time
	Assignee    *string
	Labels      []string
}

type TaskUpdate struct {
	Title       *string
	Description **string
	Priority    **uint8
	DueDate     **time.Time
	Assignee    **string
	Labels      *[]string
}

func CreateTask(task TaskCreate) (TaskID, error) {
	var taskId, err = CreateTasks([]TaskCreate{task})
	if err != nil {
		return 0, err
	}

	if len(taskId) == 0 {
		return 0, sql.ErrNoRows
	}

	return taskId[0], nil
}

func CreateTasks(tasks []TaskCreate) ([]TaskID, error) {
	if len(tasks) == 0 {
		return []TaskID{}, nil
	}

	// Build query with placeholders
	query := `
		INSERT INTO tasks (title, description, priority, due_date, assignee, labels)
		VALUES
	`
	args := []interface{}{}
	placeholderIndex := 1

	for i, task := range tasks {
		if i > 0 {
			query += ",\n"
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)",
			placeholderIndex,
			placeholderIndex+1,
			placeholderIndex+2,
			placeholderIndex+3,
			placeholderIndex+4,
			placeholderIndex+5,
		)
		args = append(args,
			task.Title,
			task.Description,
			task.Priority,
			task.DueDate,
			task.Assignee,
			pq.StringArray(task.Labels),
		)
		placeholderIndex += 6
	}

	query += "\nRETURNING id"

	// Run the query
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []TaskID
	for rows.Next() {
		var id TaskID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func GetTask(id TaskID) (*Task, error) {
	var tasks, err = GetTasks([]TaskID{id})

	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, sql.ErrNoRows
	}

	return &tasks[0], err
}

func GetTasks(ids []TaskID) ([]Task, error) {

	if len(ids) == 0 {
		return []Task{}, nil
	}

	var builder strings.Builder

	builder.WriteString(
		`
		select id, title, description, priority, due_date, assignee, labels
		from tasks
		where id in (
		`,
	)

	for i, id := range ids {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("'")
		builder.WriteString(id.String())
		builder.WriteString("'")
	}

	builder.WriteString(")")

	var rows *sql.Rows
	var err error
	rows, err = DB.Query(builder.String())

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task = make([]Task, 0, len(ids))
	for rows.Next() {
		var task Task
		rows.Scan(&task.ID, &task.Title, &task.Description, &task.Priority, &task.DueDate, &task.Assignee, &task.Labels)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func UpdateTask(id TaskID, update TaskUpdate) error {
	return UpdateTasks([]TaskID{id}, update)
}

func UpdateTasks(ids []TaskID, update TaskUpdate) error {

	if len(ids) == 0 {
		return nil
	}

	var builder strings.Builder
	builder.WriteString("update tasks set ")

	var args []interface{}
	var argIndex int = 1
	for field, value := range map[string]interface{}{
		"title":       update.Title,
		"description": update.Description,
		"priority":    update.Priority,
		"due_date":    update.DueDate,
		"assignee":    update.Assignee,
		"labels":      update.Labels,
	} {
		if value == nil || reflect.ValueOf(value).IsNil() {
			continue
		}

		if argIndex > 1 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	builder.WriteString(" where id in (")
	for i, id := range ids {
		if i > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString("'")
		builder.WriteString(id.String())
		builder.WriteString("'")
	}
	builder.WriteString(")")

	var result, err = DB.Exec(builder.String(), args...)
	print("result: ")
	print(result)

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return err
}

func DeleteTask(id TaskID) error {
	return DeleteTasks([]TaskID{id})
}

func DeleteTasks(ids []TaskID) error {

	if len(ids) == 0 {
		return nil
	}

	var builder strings.Builder
	builder.WriteString(
		`
		delete from tasks
		where id in (
		`,
	)

	for i, id := range ids {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("'")
		builder.WriteString(id.String())
		builder.WriteString("'")
	}

	builder.WriteString(")")

	_, err := DB.Exec(builder.String())

	if err != nil {
		return err
	}

	return nil
}
