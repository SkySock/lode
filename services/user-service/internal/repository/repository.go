package repository

import (
	"context"
	"fmt"

	"github.com/SkySock/lode/services/user-service/internal/db"
	"github.com/SkySock/lode/services/user-service/internal/entity"
	"github.com/google/uuid"
)

var (
	ErrDuplicate = fmt.Errorf("duplicate entry")
	ErrNotFound  = fmt.Errorf("not found")
)

type AccountRepository interface {
	GetById(ctx context.Context, qe db.QueryExecutor, id uuid.UUID) (*entity.Account, error)
	GetByUsername(ctx context.Context, qe db.QueryExecutor, username string) (*entity.Account, error)
	GetByEmail(ctx context.Context, qe db.QueryExecutor, email string) (*entity.Account, error)
	Create(ctx context.Context, qe db.QueryExecutor, account *entity.Account) (uuid.UUID, error)
}

type ProfileRepository interface {
	Create(ctx context.Context, qe db.QueryExecutor, profile *entity.Profile) (uuid.UUID, error)
	GetByID(ctx context.Context, qe db.QueryExecutor, id uuid.UUID) (*entity.Profile, error)
	GetByProfileName(ctx context.Context, qe db.QueryExecutor, name string) (*entity.Profile, error)
}

type SessionRepository interface {
	Save(ctx context.Context, token string, data *entity.Session) error
	Get(ctx context.Context, token string) (*entity.Session, error)
	Revoke(ctx context.Context, token string) error
	Delete(ctx context.Context, token string) error
	GetUserSessions(ctx context.Context, userID string) ([]string, error)
}
