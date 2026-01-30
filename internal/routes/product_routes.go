package routes

import (
	"APIGolang/internal/controller"
	"APIGolang/internal/repository"
	"APIGolang/internal/usecase"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine, db *sql.DB) {

	productRepository := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUseCase(productRepository)
	productController := controller.NewProductController(productUsecase)

	r.GET("/products", productController.GetProducts)
	r.GET("/product/:id", productController.GetProductById)
	r.POST("/product", productController.CreateProduct)
	r.PUT("/product/:id", productController.UpdateProductById)
	r.DELETE("/product/:id", productController.DeleteProductById)

}