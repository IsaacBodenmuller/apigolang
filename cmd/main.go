package main

import (
	"APIGolang/controller"
	"APIGolang/db"
	"APIGolang/repository"
	"APIGolang/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	
	server := gin.Default()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	//Camada Repository
	ProductRepository := repository.NewProductRepository(dbConnection)

	//Camada Usecase
	ProductUsecase := usecase.NewProductUseCase(ProductRepository)

	//Camada Controller
	ProductController := controller.NewProductController(ProductUsecase)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	server.GET("/products", ProductController.GetProducts)
	server.GET("/product/:id", ProductController.GetProductById)
	server.POST("/product", ProductController.CreateProduct)
	// server.PUT("/product/:id", ProductController)
	// server.PUT("/product/:id", ProductController)

	server.Run(":8000")
}