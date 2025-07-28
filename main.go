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
	r.PATCH("/tasks", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)

	r.POST("/users", handlers.CreateUser)
	r.GET("/users/:id", handlers.GetUser)
	r.PATCH("/users", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	r.POST("/lists", handlers.CreateList)
	r.GET("/lists/:id", handlers.GetList)
	r.PATCH("/lists", handlers.UpdateList)
	r.DELETE("/lists/:id", handlers.DeleteList)

	return r
}

func main() {

	db.Init()

	r := setupRouter()
	r.Run(":8080")
}
