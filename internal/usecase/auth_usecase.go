package usecase

import (
	"APIGolang/internal/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetToken(nome_usuario string) (*model.User, string, error)
	CreateUser(user model.User) error
	UserExists(nomeUsuario string) (bool, error)
	EmailExists(email string) (bool, error)
}

type AuthUseCase struct {
	repository UserRepository
}

func NewAuthUseCase(r UserRepository) *AuthUseCase {
	return &AuthUseCase{repository: r}
}

func (a *AuthUseCase) Login(nome_usuario, senha string) (*model.User, error) {
	
	user, senha_usuario, err := a.repository.GetToken(nome_usuario)
	if err != nil {
		return nil, err
	}

	if !user.Ativo {
		return nil, errors.New("O usuário está inativo")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(senha_usuario),
		[]byte(senha),
	)

	if err != nil {
		return nil, errors.New("Senha inválida")
	}

	return user, nil
}

func (a *AuthUseCase) CreateUser(req model.CreateUserRequest) error {

	userExists, err := a.repository.UserExists(req.NomeUsuario)
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

	if req.Perfil == "" {
		req.Perfil = "SEM-P"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("erro ao criptografar senha")
	}

	user := model.User{
		Nome: req.Nome,
		NomeUsuario: req.NomeUsuario,
		Email: req.Email,
		Senha: string(hash),
		Perfil: req.Perfil,
		Ativo: true,
	}

	return a.repository.CreateUser(user)
}