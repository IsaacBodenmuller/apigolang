package controller

import (
	"APIGolang/internal/model"
	"APIGolang/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	usecase *usecase.UserUseCase
}

func NewUserController(uc *usecase.UserUseCase) *UserController {
	return &UserController{usecase: uc}
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
// @Router /user/create [post]
func (userCtrl *UserController) CreateUser(c *gin.Context) {

	var req model.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			errors := make(map[string]string)

			for _, filedErr := range validationErrors {
				switch filedErr.Field() {
				case "Nome":
					errors["name"] = "Nome é obrigatório"
				case "NomeUsuario":
					errors["username"] = "Nome de usuário é obrigatório"
				case "Email":
					errors["email"] = "Email é obrigatório"
				case "Senha":
					errors["password"] = "Senha é obrigatório"
				}
			}

			c.JSON(400, gin.H{
				"errors": err.Error() + " aqui ",
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


// //@Param credentials body model.CreateUserRequest true "Dados do usuário"

// GetAllUsers godoc
// @Summary Listar todos usuário
// @Description Lista todos os usuários que existem
// @Tags User
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /user/getall [post]
func (userCtrl *UserController) GetAllUsers(c *gin.Context) {

	users, err := userCtrl.usecase.GetAllUsers()
	if err != nil {
		c.JSON(400, err)
	}

	c.JSON(200, users)

}