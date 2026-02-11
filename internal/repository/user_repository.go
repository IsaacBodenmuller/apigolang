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

func (r *UserRepository) GetToken(request_name string) (*model.User, string, error) {

	var user model.User
	var user_password string

	query := "SELECT id_usuario, nome, nome_usuario, senha, perfil, ativo FROM usuario WHERE nome_usuario = $1"

	err := r.connection.QueryRow(query, request_name).Scan(&user.Id, &user.Name, &user.Username, &user_password, &user.Profile, &user.Active)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New("Credenciais inválidas")
		}
		fmt.Println(err)
		return nil, "", err
	}

	return &user, user_password, nil
}

func (r *UserRepository) CreateUser(user model.User) error {

	query := "INSERT INTO usuario (nome, nome_usuario, email, senha, perfil, ativo)" +
			 " VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.connection.Exec(query, user.Name, user.Username, user.Email, user.Password, user.Profile, user.Active)

	return err
}

func (r *UserRepository) UserExists(user_username string) (bool, error) {

	var count int

	query := "SELECT Count(1) FROM usuario WHERE nome_usuario = $1"
	err := r.connection.QueryRow(query, user_username).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) EmailExists(user_email string) (bool, error) {

	var count int

	query := "SELECT Count(1) FROM usuario WHERE email = $1"
	err := r.connection.QueryRow(query, user_email).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil
}