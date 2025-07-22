package main

import (
	"tasks/db"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {

	r := gin.Default()

	return r
}

func main() {

	db.Init()

	r := setupRouter()
	r.Run(":8080")
}
