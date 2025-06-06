package usecase

import (
	"context"

	"github.com/google/uuid"
)

type AuthUsecase interface {
	RegisterUser(ctx context.Context, userData RegistrationInfo) (uuid.UUID, error)
	Login(ctx context.Context, login, password string) (*AuthTokens, error)
	Logout(ctx context.Context, refreshToken string) error
}

type RegistrationInfo struct {
	Username string
	Email    string
	Password string
}

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}
