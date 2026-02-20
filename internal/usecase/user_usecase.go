package usecase

import (
	"APIGolang/internal/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(user model.User) error
	UserExists(user_username string) (bool, error)
	EmailExists(user_email string) (bool, error)
	GetUserById(user_id int) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	DeleteUserById(user_id int) (bool, error)
}

type UserUseCase struct {
	repository UserRepository
}

func NewUserUseCase(r UserRepository) *UserUseCase {
	return &UserUseCase{repository: r}
}

func (uc *UserUseCase) GetUserById(id int) (*model.User, error) {
	return uc.repository.GetUserById(id)
}

func (a *UserUseCase) CreateUser(req model.CreateUserRequest) error {

	userExists, err := a.repository.UserExists(req.Username)
	if err != nil {
		return err
	}
	if userExists {
		return errors.New("Nome de usuário já cadastrado")
	}

	emailExists, err := a.repository.EmailExists(req.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return errors.New("Esse email já está cadastrado")
	}

	if req.Profile == "" {
		req.Profile = "OPERADOR"
	}
	if req.Role == "" {
		req.Role = "NO-ROLE"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("erro ao criptografar senha")
	}

	user := model.User{
		Name: req.Name,
		Username: req.Username,
		Email: req.Email,
		Password: string(hash),
		Profile: req.Profile,
		Role: req.Role,
		Active: true,
	}

	return a.repository.CreateUser(user)
}

func (a *UserUseCase) GetAllUsers() ([]model.User, error) {

	return a.repository.GetAllUsers()
}

func (a *UserUseCase) DeleteUserById(user_id int) (bool, error) {

	isSucess, err := a.repository.DeleteUserById(user_id)
	if err != nil {
		return false, err
	}
	return isSucess, nil
}