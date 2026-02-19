package usecase

import (
	"APIGolang/internal/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetToken(request_name string) (*model.User, string, error)
}

type AuthUseCase struct {
	authRepo AuthRepository
}

func NewAuthUseCase(r AuthRepository) *AuthUseCase {
	return &AuthUseCase{authRepo: r}
}

func (a *AuthUseCase) Login(request_name, request_password string) (*model.User, error) {
	
	user, user_password, err := a.authRepo.GetToken(request_name)
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