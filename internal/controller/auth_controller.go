package controller

import (
	"APIGolang/internal/auth"
	"APIGolang/internal/model"
	"APIGolang/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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
		c.JSON(400,gin.H{"error": "dados inválidos"})
		return
	}

	user, err := userCtrl.usecase.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "credenciais inválidas"})
		return
	}

	accessToken, err := auth.GenerateToken(user.Id, user.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "erro ao gerar access token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.Id)
	if err != nil {
		c.JSON(500, gin.H{"error": "erro ao gerar refresh token"})
		return
	}

	maxAge := 0
	if req.Remember != nil && *req.Remember {
		maxAge = 7 * 24 * 60 * 60
	}
	c.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		"",
		false,
		true,
	)


	c.JSON(200, gin.H{
		"access_token": accessToken,
	})
}

// Refresh godoc
// @Summary Reautenticar usuário
// @Description Obtém um novo token ao usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body model.CreateUserRequest true "Dados do usuário"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (userCtrl *UserController) Refresh(c *gin.Context) {
	
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(401, gin.H{"error": "refresh token não encontrado"})
		return
	}

	token, err := auth.ValidateToken(refreshToken)
	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "refresh token inválido"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, gin.H{"error": "claims inválidas"})
		return
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		c.JSON(401, gin.H{"error": "id inválido"})
		return
	}

	userId := int(idFloat)

	newAccessToken, _ := auth.GenerateToken(userId, "")

	c.JSON(200, gin.H{
		"access_token": newAccessToken,
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