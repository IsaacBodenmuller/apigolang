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

	err := r.connection.QueryRow(query, nome_usuario).Scan(&user.Id, &user.Nome, &user.NomeUsuario, &senha_usuario, &user.Perfil, &user.Ativo)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New("Credenciais inválidas")
		}
		fmt.Println(err)
		return nil, "", err
	}

	return &user, senha_usuario, nil
}

func (r *UserRepository) CreateUser(user model.User) error {

	query := "INSERT INTO usuario (nome, nome_usuario, email, senha, perfil, ativo)" +
			 " VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.connection.Exec(query, user.Nome, user.NomeUsuario, user.Email, user.Senha, user.Perfil, user.Ativo)

	return err
}

func (r *UserRepository) UserExists(nomeUsuario string) (bool, error) {

	var count int

	query := "SELECT Count(1) FROM usuario WHERE nome_usuario = $1"
	err := r.connection.QueryRow(query, nomeUsuario).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {

	var count int

	query := "SELECT Count(1) FROM usuario WHERE email = $1"
	err := r.connection.QueryRow(query, email).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}