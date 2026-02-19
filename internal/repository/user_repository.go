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

func (r *UserRepository) GetUserById(id int) (*model.User, error) {

	var user model.User

	query := "SELECT id_usuario, nome, nome_usuario, email, perfil, ativo FROM usuario WHERE id_usuario = $1"

	err := r.connection.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user.Profile, &user.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetToken(request_name string) (*model.User, string, error) {

	var user model.User
	var user_password string

	query := "SELECT id_usuario, nome, nome_usuario, email, senha, perfil, role, ativo FROM usuario WHERE nome_usuario = $1"

	err := r.connection.QueryRow(query, request_name).Scan(&user.Id, &user.Name, &user.Username, &user.Email, &user_password, &user.Profile, &user.Role, &user.Active)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New("Ocorreu um erro")
		}
		fmt.Println(err)
		return nil, "", err
	}

	return &user, user_password, nil
}

func (r *UserRepository) CreateUser(user model.User) error {

	query := "INSERT INTO usuario (nome, nome_usuario, email, senha, perfil, role, ativo)" +
			 " VALUES ($1, $2, $3, $4, $5, $6, $7)"

	_, err := r.connection.Exec(query, user.Name, user.Username, user.Email, user.Password, user.Profile, user.Role, user.Active)

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

func (r *UserRepository) GetAllUsers() ([]model.User, error) {

	var users []model.User

	query := "SELECT id_usuario, nome, nome_usuario, email, senha, perfil, role FROM usuario"
	rows, err := r.connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		var password string
		
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Username,
			&user.Email,
			password,
			&user.Profile,
			&user.Role,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}