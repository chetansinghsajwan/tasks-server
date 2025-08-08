package handlers

import (
	"context"
	"net/http"
	"tasks/store"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CreateTaskRequestTimeout = time.Second * 5
	GetTaskRequestTimeout    = time.Second * 5
	UpdateTaskRequestTimeout = time.Second * 5
	DeleteTaskRequestTimeout = time.Second * 5
)

type CreateTaskRequest struct {
	ListID      uint64     `json:"list_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Priority    *uint32    `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	Assignee    *string    `json:"assignee"`
	Labels      []string   `json:"labels"`
}

type UpdateTaskRequest struct {
	ListID      **uint64    `json:"list_id"`
	Title       **string    `json:"title"`
	Description **string    `json:"description"`
	Priority    **uint32    `json:"priority"`
	DueDate     **time.Time `json:"due_date"`
	Assignee    **string    `json:"assignee"`
	Labels      *[]string   `json:"labels"`
}

func CreateTask(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), CreateTaskRequestTimeout)

	defer cancel()

	var body CreateTaskRequest
	var err = c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	var taskID store.TaskID
	var serr *store.StoreError
	taskID, serr = ST.CreateTask(ctx, store.CreateTaskParams{
		ListID:      body.ListID,
		Title:       body.Title,
		Description: body.Description,
		Priority:    body.Priority,
		DueDate:     body.DueDate,
		Assignee:    body.Assignee,
		Labels:      body.Labels,
	})

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"taskID": taskID,
	})
}

func GetTask(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), GetTaskRequestTimeout)

	defer cancel()

	var taskID, err = store.ParseTaskID(c.Param("id"))
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var task *store.Task
	var serr *store.StoreError
	task, serr = ST.GetTask(ctx, taskID)

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": *task,
	})
}

func UpdateTask(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), UpdateTaskRequestTimeout)

	defer cancel()

	var taskID store.TaskID
	var err error
	if taskID, err = store.ParseTaskID(c.Param("id")); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body UpdateTaskRequest
	if err = c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	var serr *store.StoreError = ST.UpdateTask(ctx, taskID, store.UpdateTaskParams{
		ListID:      body.ListID,
		Title:       body.Title,
		Description: body.Description,
		Priority:    body.Priority,
		DueDate:     body.DueDate,
		Assignee:    body.Assignee,
		Labels:      body.Labels,
	})

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
	})
}

func DeleteTask(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(
		context.Background(), DeleteTaskRequestTimeout)

	defer cancel()

	var taskID store.TaskID
	var err error
	if taskID, err = store.ParseTaskID(c.Param("id")); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var serr *store.StoreError = ST.DeleteTask(ctx, taskID)

	if serr != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": serr.Msg,
		})
		return
	}

	c.Status(http.StatusOK)
}
