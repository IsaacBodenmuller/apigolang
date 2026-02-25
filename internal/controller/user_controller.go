package controller

import (
	"APIGolang/internal/model"
	"APIGolang/internal/usecase"
	"strconv"

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

		// if req.Role != "ADM" {
		// 	c.JSON(403, gin.H{"error": "acesso negado"})
		// }

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

	
	// if req.Role != "ADM" {
	// 	c.JSON(403, gin.H{"error": "acesso negado"})
	// }

	users, err := userCtrl.usecase.GetAllUsers()
	if err != nil {
		c.JSON(400, err)
	}

	c.JSON(200, users)

}

// func (userCtrl *UserController) GetUserById(c *gin.Context) {

// 	id := c.Param("id")
// 	user
// }

func (userCtrl *UserController) DeleteUserById(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		response := model.Response{
			Message: "Id do usuário não pode ser nulo",
		}
		c.JSON(400, response)
		return
	}
	userId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um número",
		}
		c.JSON(400, response)
		return
	}

	isSucess, err := userCtrl.usecase.DeleteUserById(userId)
	if err != nil {
		response := model.Response{
			Message: "Houve um erro",
		}
		c.JSON(400, response)
		return
	}

	c.JSON(200, isSucess)
}

func (userCtrl *UserController) UpdateUserById(c *gin.Context) {

	var user model.UpdateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Dados inválidos"})
		return
	}

	id := c.Param("id")
	if id == "" {
		response := model.Response{
			Message: "Id do usuário não pode ser nulo",
		}
		c.JSON(400, response)
		return
	}
	userId, err := strconv.Atoi(id)
	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um número",
		}
		c.JSON(400, response)
		return
	}

	isSucess, err := userCtrl.usecase.UpdateUserById(user, userId)
	if err != nil {
		c.JSON(500, err)
		return
	}
	if !isSucess {
	c.JSON(404, gin.H{"error": "Usuário não encontrado"})
	return
}

	c.JSON(200, isSucess)
}