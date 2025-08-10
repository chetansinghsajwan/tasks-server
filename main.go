package main

import (
	"tasks/db"
	"tasks/handlers"
	"tasks/services"
	"tasks/store/pg"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	var r = gin.Default()

	r.POST("/auth/login", handlers.Login)
	r.POST("/api/users", handlers.CreateUser)

	r.Group("/api", handlers.AuthenticateMiddleware)
	{
		r.POST("/tasks", handlers.CreateTask)
		r.GET("/tasks/:id", handlers.GetTask)
		r.PATCH("/tasks/:id", handlers.UpdateTask)
		r.DELETE("/tasks/:id", handlers.DeleteTask)

		r.GET("/users/:id", handlers.GetUser)
		r.PATCH("/users/:id", handlers.UpdateUser)
		r.DELETE("/users/:id", handlers.DeleteUser)

		r.POST("/lists", handlers.CreateList)
		r.GET("/lists/:id", handlers.GetList)
		r.PATCH("/lists/:id", handlers.UpdateList)
		r.DELETE("/lists/:id", handlers.DeleteList)

		r.POST("/lists/access/:id", handlers.AddListAccess)
		r.GET("/lists/access/:id", handlers.GetListAccess)
		r.DELETE("/lists/access/:id", handlers.RemoveListAccess)
	}

	return r
}

func main() {

	db.Init()

	services.ST = pg.PostgresStore{
		Pool: db.Pool,
	}

	var r = setupRouter()
	r.Run(":8080")
}
