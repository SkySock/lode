package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/SkySock/lode/services/user-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type accountRepository struct {
	dbPool *pgxpool.Pool
}

func NewAccountRepository(dbPool *pgxpool.Pool) AccountRepository {
	return &accountRepository{
		dbPool: dbPool,
	}
}

var _ AccountRepository = &accountRepository{}

func (r *accountRepository) Create(ctx context.Context, account entity.Account) (uuid.UUID, error) {
	query := `
		INSERT INTO account (id, username, email, password_hash)
			VALUES ($1, $2, $3, $4)
			RETURNING id
	`
	newId, err := uuid.NewV7()
	if err != nil {
		return uuid.UUID{}, err
	}

	var id uuid.UUID
	err = r.dbPool.QueryRow(ctx, query, newId, account.Username, account.Email, account.PasswordHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.UUID{}, ErrDuplicate
		}
		return uuid.UUID{}, fmt.Errorf("repo: create user failed: %w", err)
	}

	return id, nil
}

func (r *accountRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Account, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM account WHERE id = $1`
	var account entity.Account

	err := r.dbPool.QueryRow(ctx, query, id).Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHash, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get user failed: %w", err)
	}

	return &account, nil
}
