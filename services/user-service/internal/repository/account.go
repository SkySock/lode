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

type accountRepository struct{}

func NewAccountRepository() AccountRepository {
	return &accountRepository{}
}

func (r *accountRepository) Create(ctx context.Context, qe db.QueryExecutor, account *entity.Account) (uuid.UUID, error) {
	query := `
		INSERT INTO account (id, username, email, password_hash)
			VALUES ($1, $2, $3, $4)
	`

	newId, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	if _, err = qe.Exec(
		ctx,
		query,
		newId,
		account.Username,
		account.Email,
		account.PasswordHash,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, ErrDuplicate
		}
		return uuid.Nil, fmt.Errorf("repo: create user failed: %w", err)
	}

	return newId, nil
}

func (r *accountRepository) GetById(ctx context.Context, qe db.QueryExecutor, id uuid.UUID) (*entity.Account, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM account WHERE id = $1`
	var account entity.Account

	err := qe.QueryRow(ctx, query, id).Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHash, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get user failed: %w", err)
	}

	return &account, nil
}

func (r *accountRepository) GetByUsername(ctx context.Context, qe db.QueryExecutor, username string) (*entity.Account, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM account WHERE username = $1`
	var account entity.Account

	err := qe.QueryRow(ctx, query, username).Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHash, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get user failed: %w", err)
	}

	return &account, nil
}

func (r *accountRepository) GetByEmail(ctx context.Context, qe db.QueryExecutor, email string) (*entity.Account, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM account WHERE email = $1`
	var account entity.Account

	err := qe.QueryRow(ctx, query, email).Scan(&account.ID, &account.Username, &account.Email, &account.PasswordHash, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get user failed: %w", err)
	}
	return &account, nil
}
