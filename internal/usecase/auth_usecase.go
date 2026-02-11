package usecase

import (
	"APIGolang/internal/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetToken(request_name string) (*model.User, string, error)
	CreateUser(user model.User) error
	UserExists(user_username string) (bool, error)
	EmailExists(email string) (bool, error)
}

type AuthUseCase struct {
	repository UserRepository
}

func NewAuthUseCase(r UserRepository) *AuthUseCase {
	return &AuthUseCase{repository: r}
}

func (a *AuthUseCase) Login(request_name, request_password string) (*model.User, error) {
	
	user, user_password, err := a.repository.GetToken(request_name)
	if err != nil {
		return nil, err
	}

	if !user.Active {
		return nil, errors.New("O usuário está inativo")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user_password),
		[]byte(request_password),
	)

	if err != nil {
		return nil, errors.New("Senha inválida")
	}

	return user, nil
}

func (a *AuthUseCase) CreateUser(req model.CreateUserRequest) error {

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
		req.Profile = "SEM-P"
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
		Active: true,
	}

	return a.repository.CreateUser(user)
}