package controller

import (
	"APIGolang/internal/auth"
	"APIGolang/internal/model"
	"APIGolang/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	authUsecase *usecase.AuthUseCase
	userUsecase *usecase.UserUseCase
}

func NewAuthController(authUc *usecase.AuthUseCase, userUc *usecase.UserUseCase) *AuthController {
	return &AuthController{authUsecase: authUc, userUsecase: userUc}
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
func (authCtrl *AuthController) Login(c *gin.Context) {

	var req model.TokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if req.Remember == nil {
			c.JSON(400,gin.H{"error": "informe a tag remember"})
			return
		}
		c.JSON(400,gin.H{"error": "dados inválidos"})
		return
	}

	user, err := authCtrl.authUsecase.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "credenciais inválidas"})
		return
	}

	accessToken, err := auth.GenerateToken(user.Id, user.Username, user.Email, user.Profile, user.Role, user.Active)
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
		maxAge = 60 * 60
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

// //@Param credentials body model.CreateUserRequest true "Dados do usuário"

// Refresh godoc
// @Summary Reautenticar usuário
// @Description Obtém um novo token ao usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (authCtrl *AuthController) Refresh(c *gin.Context) {
	
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

	claims := token.Claims.(jwt.MapClaims)
	idFloat := claims["id"].(float64)
	userId := int(idFloat)

	user, err := authCtrl.userUsecase.GetUserById(userId)
	if err != nil {
		c.JSON(401, gin.H{"error": "usuário não encontrado"})
		return
	}

	newAccessToken, err := auth.GenerateToken(user.Id, user.Username, user.Email, user.Profile, user.Role, user.Active)
	if err != nil {
		c.JSON(500, gin.H{"error": "erro ao gerar access token"})
		return
	}
	
	c.JSON(200, gin.H{
		"access_token": newAccessToken,
	})
}


// //@Param credentials body model.CreateUserRequest true "Dados do usuário"

// AlterPassword godoc
// @Summary Alterar senha
// @Description Altera a senha de um usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /auth/alterpassword [post]
func (authCtrl *AuthController) AlterPassword(c *gin.Context) {

}
// 	var req model.CreateUserRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {

// 		validationErrors, ok := err.(validator.ValidationErrors)
// 		if ok {
// 			errors := make(map[string]string)

// 			for _, filedErr := range validationErrors {
// 				switch filedErr.Field() {
// 				case "Nome":
// 					errors["name"] = "Nome é obrigatório"
// 				case "NomeUsuario":
// 					errors["username"] = "Nome de usuário é obrigatório"
// 				case "Email":
// 					errors["email"] = "Email é obrigatório"
// 				case "Senha":
// 					errors["password"] = "Senha é obrigatório"
// 				}
// 			}

// 			c.JSON(400, gin.H{
// 				"errors": err.Error() + " aqui ",
// 			})
// 			return
// 		}
// 		c.JSON(400, gin.H{"error": "Dados inválidos"})
// 		return
// 	}

// 	err := authCtrl.usecase.CreateUser(req)
// 	if err != nil {
// 		c.JSON(400, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(201, gin.H{
// 		"message": "Usuário criado com sucesso",
// 	})
// }


