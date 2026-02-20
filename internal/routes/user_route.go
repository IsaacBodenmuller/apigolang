package routes

import (
	"APIGolang/internal/controller"
	"APIGolang/internal/middleware"
	"APIGolang/internal/repository"
	"APIGolang/internal/usecase"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, db *sql.DB) {
	
	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUseCase(&userRepository)
	userController := controller.NewUserController(userUsecase)
	userRoutes := r.Group("/user")

	userRoutes.Use(middleware.JWTAuth()) 
	{
		userRoutes.POST("/create", userController.CreateUser)
		userRoutes.GET("/getAll", userController.GetAllUsers)
		userRoutes.DELETE("/delete/:id", userController.DeleteUserById)
	}
}

