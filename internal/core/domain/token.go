package domain

import "time"

type RefreshToken struct {
	ID        int
	IDUser    int
	Token     string
	ExpiresAt time.Time
	CreatedAt *time.Time
}
