package main

import (
	"tasks/db"
	"tasks/routes"
	"tasks/services"
	"tasks/store/pg"
)

func main() {

	db.Init()

	services.ST = pg.PostgresStore{
		Pool: db.Pool,
	}

	var r = routes.SetupRouter()
	r.Run(":8080")
}
