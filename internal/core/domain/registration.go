package domain

import (
	"time"

	"github.com/google/uuid"
)

type CSVRow struct {
	FIO         string
	Email       string
	Role        Role
	PhoneNumber *string
	GroupName   string
}

type RegistrationRequest struct {
	ID          uuid.UUID
	FIO         string
	Role        Role
	PhoneNumber *string
	Email       string
	GroupName   string
	Token       uuid.UUID
	Status      string
	ExpiresAt   time.Time
}

type UploadResult struct {
	Total     int          `json:"total"`
	Processed int          `json:"processed"`
	Skipped   int          `json:"skipped"`
	Warnings  []WarningRow `json:"warnings"`
}

type WarningRow struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
}
