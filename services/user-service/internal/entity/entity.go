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
	JoinedAt    time.Time
}
