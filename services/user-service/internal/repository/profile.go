package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/SkySock/lode/services/user-service/internal/db"
	"github.com/SkySock/lode/services/user-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type profileRepository struct{}

func NewProfileRepository() ProfileRepository {
	return &profileRepository{}
}

func (r *profileRepository) Create(ctx context.Context, qe db.QueryExecutor, profile *entity.Profile) (uuid.UUID, error) {
	query := `
		INSERT INTO profile(id, user_id, profile_name, display_name, bio, avatar) 
			VALUES ($1, $2, $3, $4, $5, $6)
	`

	newId, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	if _, err = qe.Exec(
		ctx,
		query,
		newId,
		profile.UserID,
		profile.ProfileName,
		profile.DisplayName,
		profile.Bio,
		profile.Avatar,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, ErrDuplicate
		}
		return uuid.Nil, fmt.Errorf("repo: create user profile failed: %w", err)
	}

	return newId, nil
}

func (r *profileRepository) GetByID(ctx context.Context, qe db.QueryExecutor, id uuid.UUID) (*entity.Profile, error) {
	query := `
		SELECT id, user_id, profile_name, display_name, bio, avatar, created_at
			FROM profile 
			WHERE id = $1
	`

	var profile entity.Profile
	err := qe.QueryRow(ctx, query, id).
		Scan(
			&profile.ID,
			&profile.UserID,
			&profile.ProfileName,
			&profile.DisplayName,
			&profile.Bio,
			&profile.Avatar,
			&profile.CreatedAt,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get user profile failed: %w", err)
	}

	return &profile, nil
}

func (r *profileRepository) GetByProfileName(ctx context.Context, qe db.QueryExecutor, name string) (*entity.Profile, error) {
	query := `
		SELECT id, user_id, profile_name, display_name, bio, avatar, created_at
			FROM profile 
			WHERE profile_name = $1
	`

	var profile entity.Profile

	if err := qe.QueryRow(ctx, query, name).
		Scan(
			&profile.ID,
			&profile.UserID,
			&profile.ProfileName,
			&profile.DisplayName,
			&profile.Bio,
			&profile.Avatar,
			&profile.CreatedAt,
		); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get user profile failed: %w", err)
	}

	return &profile, nil
}
