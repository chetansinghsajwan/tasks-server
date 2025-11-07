package main

import (
	"os"
	"fmt"

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

	var port = os.Getenv("SERVER_PORT")
	if port == "" {
		panic("SERVER_PORT environment variable is not set")
	}

	r.Run(fmt.Sprintf(":%s", port))
}
