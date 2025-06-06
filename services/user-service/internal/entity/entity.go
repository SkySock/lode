package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type Profile struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	ProfileName string
	DisplayName string
	Bio         string
	Avatar      string
	CreatedAt   time.Time
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProfileID string    `json:"profile_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}
