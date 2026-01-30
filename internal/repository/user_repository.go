package repository

import (
	"APIGolang/internal/model"
	"database/sql"
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

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {

	var user model.User

	query := "SELECT user_id, user_email, user_password FROM user WHERE user_email = $1"

	err := r.connection.QueryRow(query, email).Scan(&user.Id, &user.Email, &user.Password)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}