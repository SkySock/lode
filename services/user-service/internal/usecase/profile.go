package usecase

import (
	"context"

	"github.com/SkySock/lode/services/user-service/internal/entity"
	"github.com/SkySock/lode/services/user-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type profileUsecase struct {
	pgPool *pgxpool.Pool
	repo   repository.ProfileRepository
}

func NewProfileUsecase(pgPool *pgxpool.Pool, repo repository.ProfileRepository) ProfileUsecase {
	return &profileUsecase{
		pgPool: pgPool,
		repo:   repo,
	}
}

func (uc *profileUsecase) GetUserProfile(ctx context.Context, id uuid.UUID) (*entity.Profile, error) {
	profile, err := uc.repo.GetByID(ctx, uc.pgPool, id)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
