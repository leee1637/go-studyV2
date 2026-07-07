package domain

import "time"

type RefreshToken struct {
	ID         int
	IDUser     int
	Token      string
	Expires_at time.Time
	Created_at *time.Time
}
