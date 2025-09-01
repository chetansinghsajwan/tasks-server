package routes

import (
	"tasks/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	var r = gin.Default()

	r.POST("/auth/login", handlers.Login)

	var api = r.Group("/api")
	{
		api.POST("/users", handlers.CreateUser)

		var auth = api.Group("", handlers.AuthenticateMiddleware)
		{
			auth.POST("/tasks", handlers.CreateTask)
			auth.GET("/tasks/:id", handlers.GetTask)
			auth.PATCH("/tasks/:id", handlers.UpdateTask)
			auth.DELETE("/tasks/:id", handlers.DeleteTask)

			auth.GET("/users/:id", handlers.GetUser)
			auth.PATCH("/users/:id", handlers.UpdateUser)
			auth.DELETE("/users/:id", handlers.DeleteUser)

			auth.POST("/lists", handlers.CreateList)
			auth.GET("/lists/:id", handlers.GetList)
			auth.PATCH("/lists/:id", handlers.UpdateList)
			auth.DELETE("/lists/:id", handlers.DeleteList)

			auth.POST("/lists/access/:id", handlers.AddListAccess)
			auth.GET("/lists/access/:id", handlers.GetListAccess)
			auth.DELETE("/lists/access/:id", handlers.RemoveListAccess)
		}
	}

	return r
}
