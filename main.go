package main

import (
	"notes-app/src/app"
	"notes-app/src/commons/database"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	database.ConnectPostgresql(".env")

	app.RegisterRouter(router, database.DB)

	router.Run(":3000")
}
