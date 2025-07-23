package main

import (
	"tasks/db"
	"tasks/handlers"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.POST("/tasks", handlers.CreateTask)
	r.GET("/tasks/:id", handlers.GetTask)
	r.PATCH("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	return r
}

func main() {

	db.Init()

	r := setupRouter()
	r.Run(":8080")
}
