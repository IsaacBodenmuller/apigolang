package controller

import (
	"APIGolang/internal/auth"
	"APIGolang/internal/model"
	"APIGolang/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	usecase *usecase.AuthUseCase
}

func NewUserController(uc *usecase.AuthUseCase) *UserController {
	return &UserController{usecase: uc}
}

// Login godoc
// @Summary Autenticar usuário
// @Description Realiza login e retorna o token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body model.TokenRequest true "Credenciais do usuário"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (userCtrl *UserController) Login(c *gin.Context) {

	var req model.TokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400,gin.H{
		"error": err.Error()})
		return
	}

	user, err := userCtrl.usecase.Login(req.NomeUsuario, req.Senha)
	if err != nil {
		c.JSON(401, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, _ := auth.GenerateToken(user.Id, user.Nome, user.NomeUsuario, user.Perfil)

	c.JSON(200, gin.H{
		"token": token,
	})
}

// CreateUser godoc
// @Summary Criar usuário
// @Description Cria um novo usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body model.CreateUserRequest true "Dados do usuário"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /auth/create [post]
func (userCtrl *UserController) CreateUser(c *gin.Context) {

	var req model.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			errors := make(map[string]string)

			for _, filedErr := range validationErrors {
				switch filedErr.Field() {
				case "Nome":
					errors["nome"] = "Nome é obrigatório"
				case "NomeUsuario":
					errors["nomeUsuario"] = "Nome de usuário é obrigatório"
				case "Email":
					errors["email"] = "Email é obrigatório"
				case "Senha":
					errors["senha"] = "Senha é obrigatório"
				}
			}

			c.JSON(400, gin.H{
				"errors": errors,
			})
			return
		}
		c.JSON(400, gin.H{"error": "Dados inválidos"})
		return
	}

	err := userCtrl.usecase.CreateUser(req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Usuário criado com sucesso",
	})
}