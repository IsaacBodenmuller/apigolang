package routes

import (
	"APIGolang/internal/controller"
	"APIGolang/internal/repository"
	"APIGolang/internal/usecase"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, db *sql.DB) {
	
	userRepository := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUseCase(&userRepository)
	authController := controller.NewUserController(authUsecase)

	r.POST("/auth/login", authController.Login)
	r.POST("/auth/create", authController.CreateUser)
}