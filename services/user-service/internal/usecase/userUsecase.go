package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/SkySock/lode/libs/utils/argon2id"
	"github.com/SkySock/lode/services/user-service/internal/entity"
	repo "github.com/SkySock/lode/services/user-service/internal/repository"
	"github.com/google/uuid"
)

var ErrEmailOrUsernameAlreadyExists = errors.New("email or username already exists")

type userUsecase struct {
	accRepo repo.AccountRepository
}

func NewUserUsecase(accRepo repo.AccountRepository) UserUsecase {
	return &userUsecase{
		accRepo: accRepo,
	}
}

func (u *userUsecase) RegisterUser(ctx context.Context, userData RegistrationInfo) (uuid.UUID, error) {
	account := entity.Account{}
	account.Email = userData.Email
	account.Username = userData.Username

	params := &argon2id.Params{
		Memory:      46 * 1024,
		Iterations:  2,
		Parallelism: 1,
		KeyLength:   64,
		SaltLength:  16,
	}

	passwordHash, err := argon2id.HashPassword([]byte(userData.Password), params)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("hashing password error: %w", err)
	}
	account.PasswordHash = passwordHash

	accountId, err := u.accRepo.Create(ctx, account)
	if err != nil {
		if errors.Is(err, repo.ErrDuplicate) {
			return uuid.UUID{}, ErrEmailOrUsernameAlreadyExists
		}
		return uuid.UUID{}, fmt.Errorf("failed to create account: %w", err)
	}

	return accountId, nil
}
