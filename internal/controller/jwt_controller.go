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

func (userCtrl *UserController) Login(c *gin.Context) {

	var req model.TokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400,gin.H{
		"error": err.Error()})
		return
	}

	user, err := userCtrl.usecase.Login(req.Nome_Usuario, req.Senha)
	if err != nil {
		c.JSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, _ := auth.GenerateToken(user.Id, user.Nome, user.Nome_Usuario, user.Perfil)

	c.JSON(200, gin.H{
		"token": token,
	})
}