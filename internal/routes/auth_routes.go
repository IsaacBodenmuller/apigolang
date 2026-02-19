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
	userUsecase := usecase.NewUserUseCase(&userRepository)
	authController := controller.NewAuthController(authUsecase, userUsecase)

	r.POST("/auth/login", authController.Login)
	r.POST("/auth/refresh", authController.Refresh)
	r.POST("/auth/alterPassword", authController.AlterPassword)

}
