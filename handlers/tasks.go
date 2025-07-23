package handlers

import (
	"net/http"
	"tasks/db"
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

	var task TaskCreate
	var err = c.BindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	var taskId db.TaskID
	taskId, err = db.CreateTask(db.TaskCreate{
		Title:       task.Title,
		Description: &task.Description,
		Priority:    &task.Priority,
		DueDate:     &task.DueDate,
		Assignee:    &task.Assignee,
		Labels:      task.Labels,
	})

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

	var taskId, err = db.ParseTaskID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var task *db.Task
	task, err = db.GetTask(taskId)

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

	var taskId, err = db.ParseTaskID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	var taskUpdate db.TaskUpdate
	err = c.BindJSON(&taskUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = db.UpdateTask(taskId, taskUpdate)
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

	var taskId, err = db.ParseTaskID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.DeleteTask(taskId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
