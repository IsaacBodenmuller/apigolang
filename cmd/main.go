package main

import (
	"APIGolang/internal/db"
	"APIGolang/internal/routes"

	"github.com/gin-gonic/gin"

	_ "APIGolang/swagger/v1"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Mercado
// @version 1.0
// @description API para o sistema do mercado

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	
	server := gin.Default()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	routes.RegisterProductRoutes(server, dbConnection)
	routes.RegisterAuthRoutes(server, dbConnection)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.Run(":8000")
}