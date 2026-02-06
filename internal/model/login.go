package model

type TokenRequest struct {
	NomeUsuario string `json:"nomeUsuario"`
	Senha       string `json:"senha"`
}