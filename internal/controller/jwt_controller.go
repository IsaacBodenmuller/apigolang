package controller

import (
	"APIGolang/internal/auth"
	"APIGolang/internal/model"
	"APIGolang/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	usecase *usecase.AuthUseCase
}

func NewUserController(uc *usecase.AuthUseCase) *UserController {
	return &UserController{usecase: uc}
}

func (uc *UserController) Login(c *gin.Context) {

	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400,gin.H{
		"error": err.Error()})
		return
	}

	user, err := uc.usecase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "credenciais inválidas"})
		return
	}

	token, _ := auth.GenerateToken(user.Id, user.Email)

	c.JSON(200, gin.H{
		"token": token,
	})
}