package main

import (
	"APIGolang/internal/db"
	"APIGolang/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	
	server := gin.Default()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	routes.RegisterProductRoutes(server, dbConnection)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.Run(":8000")
}