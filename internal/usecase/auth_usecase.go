package usecase

import (
	"APIGolang/internal/model"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetToken(request_name string) (*model.User, string, error)
	ChangePassword(email, password string) (bool, error)
	EmailExists(user_email string) (bool, error)
}

type AuthUseCase struct {
	authRepo AuthRepository
	userRepo UserRepository
}

func NewAuthUseCase(authRepo AuthRepository, userRepo UserRepository) *AuthUseCase {
	return &AuthUseCase{authRepo: authRepo, userRepo: userRepo}
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

func (a *AuthUseCase) ChangePassword(user_request model.ChangePassword) (bool, error) {

	emailExist, err := a.authRepo.EmailExists(user_request.Email)
	if err != nil {
		return false, err
	}
	if !emailExist {
		return false, errors.New("Esse email não existe")
	}

	// user, err := a.userRepo.GetUserByEmail(user_request.Email)
	// if err != nil {
	// 	return false, errors.New("Esse email não está cadastrado")
	// }

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(user_request.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return false, err
	}

	return a.authRepo.ChangePassword(user_request.Email, string(hash))
}