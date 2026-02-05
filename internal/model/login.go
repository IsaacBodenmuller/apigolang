package model

type TokenRequest struct {
	Nome_Usuario string `json:"nome_usuario"`
	Senha        string `json:"senha"`
}