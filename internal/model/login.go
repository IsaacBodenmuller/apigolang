package model

type TokenRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remember *bool  `json:"remember" binding:"required"`
}