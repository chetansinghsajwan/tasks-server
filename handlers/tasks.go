package handlers

import (
	"context"
	"net/http"
	"tasks/db"
	"tasks/sqlc"
	"tasks/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskCreate struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Priority    uint8     `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	Assignee    string    `json:"assignee"`
	Labels      []string  `json:"labels"`
}

func CreateTask(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var task sqlc.CreateTaskParams
	var err = c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	var taskId int32
	taskId, err = db.Queries.CreateTask(ctx, task)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"taskId": taskId,
	})
}

func GetTask(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var taskId, err = utils.ParseInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var task sqlc.Task
	task, err = db.Queries.GetTask(ctx, taskId)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": gin.H{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"priority":    task.Priority,
			"dueDate":     task.DueDate,
			"assignee":    task.Assignee,
			"labels":      task.Labels,
		},
	})
}

func UpdateTask(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var taskUpdate *sqlc.UpdateTaskParams
	var err error
	err = c.BindJSON(&taskUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = db.Queries.UpdateTask(ctx, *taskUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
	})
}

func DeleteTask(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var taskId, err = utils.ParseInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.Queries.DeleteTask(ctx, taskId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
