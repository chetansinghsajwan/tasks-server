package store

import (
	"time"
)

type Task struct {
	ID          uint64
	ListID      uint64
	Title       string
	Description *string
	Priority    *uint32
	DueDate     *time.Time
	Assignee    *string
	Labels      []string
}

type CreateTaskParams struct {
	ListID      uint64
	Title       string
	Description *string
	Priority    *uint32
	DueDate     *time.Time
	Assignee    *string
	Labels      []string
}

type UpdateTaskParams struct {
	ListID      **uint64
	Title       **string
	Description **string
	Priority    **uint32
	DueDate     **time.Time
	Assignee    **string
	Labels      *[]string
}
