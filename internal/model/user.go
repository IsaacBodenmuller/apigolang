package model

type User struct {
	Id          int
	Nome        string
	NomeUsuario string
	Email       string
	Senha       string
	Perfil      string
	Ativo       bool
}