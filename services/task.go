package services

import (
	"tasks/errorcodes"
	"tasks/store"
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

func CreateTask(ctx ServiceContext, args CreateTaskParams) (uint64, *ServiceError) {

	var taskID uint64
	var serr *store.StoreError
	taskID, serr = ST.CreateTask(ctx.Ctx, store.CreateTaskParams{
		ListID:      args.ListID,
		Title:       args.Title,
		Description: args.Description,
		Priority:    args.Priority,
		DueDate:     args.DueDate,
		Assignee:    args.Assignee,
		Labels:      args.Labels,
	})

	if serr != nil {

		return 0, &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return taskID, nil
}

func GetTask(ctx ServiceContext, id uint64) (*Task, *ServiceError) {

	var task *store.Task
	var serr *store.StoreError
	task, serr = ST.GetTask(ctx.Ctx, id)

	if serr != nil {

		return nil, &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return &Task{
		ID:          task.ID,
		ListID:      task.ListID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		Assignee:    task.Assignee,
		Labels:      task.Labels,
	}, nil
}

func UpdateTask(ctx ServiceContext, id uint64, args UpdateTaskParams) *ServiceError {

	var serr *store.StoreError = ST.UpdateTask(
		ctx.Ctx,
		id,
		store.UpdateTaskParams{
			ListID:      args.ListID,
			Title:       args.Title,
			Description: args.Description,
			Priority:    args.Priority,
			DueDate:     args.DueDate,
			Assignee:    args.Assignee,
			Labels:      args.Labels,
		},
	)

	if serr != nil {

		return &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return nil
}

func DeleteTask(ctx ServiceContext, id uint64) *ServiceError {

	var serr *store.StoreError = ST.DeleteTask(ctx.Ctx, id)

	if serr != nil {

		return &ServiceError{
			Code: errorcodes.Internal,
		}
	}

	return nil
}
