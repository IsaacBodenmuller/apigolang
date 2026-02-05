package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("GolangMercadoApi")

func GenerateToken(userId int, nome, nome_usuario, perfil string) (string, error) {

	claims := jwt.MapClaims{
		"id_usuario": userId,
		"nome": nome,
		"nome_usuario": nome_usuario,
		"perfil": perfil,
		"expires_in":     time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
		return secretKey, nil
	})
}