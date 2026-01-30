package middleware

import (
	"APIGolang/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token não informado",
			})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		_, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token inválido",
			})
			return
		}

		c.Next()
	}
}