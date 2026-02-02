package routes

import (
	"APIGolang/internal/controller"
	"APIGolang/internal/middleware"
	"APIGolang/internal/repository"
	"APIGolang/internal/usecase"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine, db *sql.DB) {

	productRepository := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUseCase(productRepository)
	productController := controller.NewProductController(productUsecase)
	productsRoutes := r.Group("/product")
	
	productsRoutes.Use(middleware.JWTAuth()) 
	{
		productsRoutes.GET("", productController.GetProducts)
		productsRoutes.GET("/:id", productController.GetProductById)
		productsRoutes.POST("", productController.CreateProduct)
		productsRoutes.PUT("/:id", productController.UpdateProductById)
		productsRoutes.DELETE("/:id", productController.DeleteProductById)
	}
}