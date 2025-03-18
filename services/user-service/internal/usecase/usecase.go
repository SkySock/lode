package usecase

import (
	"context"

	"github.com/google/uuid"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, userData RegistrationInfo) (uuid.UUID, error)
}

type RegistrationInfo struct {
	Username string
	Email    string
	Password string
}
