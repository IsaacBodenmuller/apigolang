package model

type CreateUserRequest struct {
	Nome        string `json:"nome" binding:"required"`
	NomeUsuario string `json:"nomeUsuario" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Senha       string `json:"senha" binding:"required"`
	Perfil      string `json:"perfil"`
}