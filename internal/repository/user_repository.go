package repository

import (
	"APIGolang/internal/model"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository struct {
	connection *sql.DB
}

func NewUserRepository(connection *sql.DB) UserRepository {
	return UserRepository{
		connection: connection,
	}
}

func (r *UserRepository) GetToken(nome_usuario string) (*model.User, string, error) {

	var user model.User
	var senha_usuario string

	query := "SELECT id_usuario, nome, nome_usuario, senha, perfil, ativo FROM usuario WHERE nome_usuario = $1"

	err := r.connection.QueryRow(query, nome_usuario).Scan(&user.Id, &user.Nome, &user.Nome_Usuario, &senha_usuario, &user.Perfil, &user.Ativo)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New("Credenciais inválidas")
		}
		fmt.Println(err)
		return nil, "", err
	}

	return &user, senha_usuario, nil
}