package usecase

import "APIGolang/internal/model"

type UserRepository interface {
	GetByEmail(email string) (*model.User, error)
}

type AuthUseCase struct {
	repository UserRepository
}

func NewAuthUseCase(r UserRepository) *AuthUseCase {
	return &AuthUseCase{repository: r}
}

func (a *AuthUseCase) Login(email, password string) (*model.User, error) {
	
	user, err := a.repository.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return  nil, err
	}

	return user, nil
}