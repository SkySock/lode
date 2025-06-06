package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SkySock/lode/libs/utils/argon2id"
	"github.com/SkySock/lode/services/user-service/internal/config"
	"github.com/SkySock/lode/services/user-service/internal/entity"
	repo "github.com/SkySock/lode/services/user-service/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailOrUsernameAlreadyExists = errors.New("email or username already exists")
	ErrIncorrectPassword            = errors.New("incorrect password")
	ErrUsernameNotFound             = errors.New("username not found")
	ErrEmailNotFound                = errors.New("email not found")
)

type authUsecase struct {
	pgPool      *pgxpool.Pool
	accountRepo repo.AccountRepository
	sessionRepo repo.SessionRepository
	profileRepo repo.ProfileRepository
	authConfig  config.AuthConfig
}

func NewAuthUsecase(
	pool *pgxpool.Pool,
	accountRepo repo.AccountRepository,
	sessionRepo repo.SessionRepository,
	profileRepo repo.ProfileRepository,
	authConfig config.AuthConfig,
) AuthUsecase {
	return &authUsecase{
		pgPool:      pool,
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
		profileRepo: profileRepo,
		authConfig:  authConfig,
	}
}

func (u *authUsecase) RegisterUser(ctx context.Context, userData RegistrationInfo) (uuid.UUID, error) {
	tx, err := u.pgPool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Begin tx: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
	}()

	account := entity.Account{}
	account.Email = strings.ToLower(userData.Email)
	account.Username = strings.ToLower(userData.Username)

	params := &argon2id.Params{
		Memory:      46 * 1024,
		Iterations:  2,
		Parallelism: 1,
		KeyLength:   64,
		SaltLength:  16,
	}

	passwordHash, err := argon2id.HashPassword([]byte(userData.Password), params)
	if err != nil {
		_ = tx.Rollback(ctx)
		return uuid.Nil, fmt.Errorf("hashing password error: %w", err)
	}
	account.PasswordHash = passwordHash

	accountId, err := u.accountRepo.Create(ctx, tx, &account)
	if err != nil {
		_ = tx.Rollback(ctx)

		if errors.Is(err, repo.ErrDuplicate) {
			return uuid.Nil, ErrEmailOrUsernameAlreadyExists
		}
		return uuid.Nil, fmt.Errorf("failed to create account: %w", err)
	}

	profile := entity.Profile{}
	profile.UserID = accountId
	profile.ProfileName = account.Username

	if _, err := u.profileRepo.Create(ctx, tx, &profile); err != nil {
		_ = tx.Rollback(ctx)

		return uuid.Nil, fmt.Errorf("failed to create profile: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("tx.Commit: %w", err)
	}

	return accountId, nil
}

func (u *authUsecase) Login(ctx context.Context, login, password string) (*AuthTokens, error) {
	account, err := u.checkUserCredentials(ctx, login, password)
	if err != nil {
		return nil, err
	}

	profile, err := u.profileRepo.GetByProfileName(ctx, u.pgPool, account.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile by username: %w", err)
	}

	accessJWT, err := u.generateAccessJWT(account.ID, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access jwt: %w", err)
	}

	refresh, err := u.generateRefreshToken(ctx, account.ID, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthTokens{
		AccessToken:  accessJWT,
		RefreshToken: refresh,
	}, nil
}

func (u *authUsecase) Logout(ctx context.Context, refreshToken string) error {
	err := u.sessionRepo.Delete(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (u *authUsecase) checkUserCredentials(ctx context.Context, login, password string) (*entity.Account, error) {
	var account *entity.Account
	var err error

	if strings.Contains(login, "@") {
		account, err = u.accountRepo.GetByEmail(ctx, u.pgPool, login)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return nil, ErrEmailNotFound
			}
			return nil, err
		}
	} else {
		account, err = u.accountRepo.GetByUsername(ctx, u.pgPool, login)
		if err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return nil, ErrUsernameNotFound
			}
			return nil, err
		}
	}

	check, err := argon2id.VerifyPassword([]byte(password), account.PasswordHash)
	if err != nil {
		return nil, err
	}

	if check {
		return account, nil
	}

	return nil, ErrIncorrectPassword
}

func (u *authUsecase) generateAccessJWT(userId, profileId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":     userId.String(),
		"profile": profileId.String(),
		"exp":     time.Now().Add(u.authConfig.Lifetime.Access).Unix(),
		"iss":     "user-service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(u.authConfig.AccessSecretKey))
}

func (u *authUsecase) generateRefreshToken(ctx context.Context, userId, profileId uuid.UUID) (string, error) {
	token, err := u.generateOpaqueToken()
	if err != nil {
		return "", err
	}
	session := entity.Session{
		ID:        token,
		UserID:    userId.String(),
		ProfileID: profileId.String(),
		ExpiresAt: time.Now().Add(u.authConfig.Lifetime.Refresh),
		Revoked:   false,
	}

	if err = u.sessionRepo.Save(ctx, token, &session); err != nil {
		return "", err
	}

	return token, nil
}

func (u *authUsecase) generateOpaqueToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
