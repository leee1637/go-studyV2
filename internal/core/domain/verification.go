package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID        uuid.UUID
	UserID    int
	Token     uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
}
