package repository

import (
	"context"
	"fmt"

	"github.com/SkySock/lode/services/user-service/internal/entity"
	"github.com/google/uuid"
)

var (
	ErrDuplicate = fmt.Errorf("duplicate entry")
	ErrNotFound  = fmt.Errorf("not found")
)

type AccountRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*entity.Account, error)
	Create(ctx context.Context, account entity.Account) (uuid.UUID, error)
}
