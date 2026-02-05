package usecase

import (
	"APIGolang/internal/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetToken(nome_usuario string) (*model.User, string, error)
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