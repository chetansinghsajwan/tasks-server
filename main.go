package main

import (
	"tasks/db"
	"tasks/handlers"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.POST("/auth/login", handlers.Login)
	r.POST("/api/users", handlers.CreateUser)

	r.Group("/api", handlers.AuthenticateMiddleware)
	{
		r.POST("/tasks", handlers.CreateTask)
		r.GET("/tasks/:id", handlers.GetTask)
		r.PATCH("/tasks", handlers.UpdateTask)
		r.DELETE("/tasks/:id", handlers.DeleteTask)

		r.GET("/users/:id", handlers.GetUser)
		r.PATCH("/users/:id", handlers.UpdateUser)
		r.DELETE("/users/:id", handlers.DeleteUser)

		r.POST("/lists", handlers.CreateList)
		r.GET("/lists/:id", handlers.GetList)
		r.PATCH("/lists", handlers.UpdateList)
		r.DELETE("/lists/:id", handlers.DeleteList)
		r.GET("/lists/access/:listid/:userid", handlers.GetListAccess)
		r.PUT("/lists/access/:listid/:userid", handlers.SetListAccess)
	}

	return r
}

func main() {

	db.Init()

	r := setupRouter()
	r.Run(":8080")
}
