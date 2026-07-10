package repository

import (
	"context"
	"errors"
	"fmt"
	"study/internal/core/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RegistrationRepository struct {
	pool *pgxpool.Pool
}

func NewRegistrationRepository(pool *pgxpool.Pool) *RegistrationRepository {
	return &RegistrationRepository{pool: pool}
}

func (r *RegistrationRepository) SaveBatch(ctx context.Context, requests []domain.RegistrationRequest) error {
	query := `INSERT INTO registration_requests 
        (id, fio, role, phone_number, email, group_name, token, expires_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	for _, req := range requests {
		_, err := r.pool.Exec(ctx, query,
			req.ID, req.FIO, string(req.Role), req.PhoneNumber, req.Email, req.GroupName, req.Token, req.ExpiresAt)
		if err != nil {
			return fmt.Errorf("ошибка сохранения заявки для %s: %w", req.Email, err)
		}
	}
	return nil
}

func (r *RegistrationRepository) GetByToken(ctx context.Context, token uuid.UUID) (domain.RegistrationRequest, error) {
	query := `SELECT id, fio, role, phone_number, email, group_name, token, status, expires_at 
              FROM registration_requests WHERE token = $1`

	var req domain.RegistrationRequest
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&req.ID, &req.FIO, &req.Role, &req.PhoneNumber, &req.Email,
		&req.GroupName, &req.Token, &req.Status, &req.ExpiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RegistrationRequest{}, fmt.Errorf("заявка не найдена")
		}
		return domain.RegistrationRequest{}, fmt.Errorf("Ошибка БД: %w", err)
	}
	return req, nil
}

func (r *RegistrationRepository) EmailExistsInPending(ctx context.Context, emails []string) ([]string, error) {
	if len(emails) == 0 {
		return []string{}, nil
	}
	query := `SELECT email FROM registration_requests 
              WHERE email = ANY($1) AND status = 'pending'`

	rows, err := r.pool.Query(ctx, query, emails)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var existing []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		existing = append(existing, email)
	}

	return existing, nil
}

func (r *RegistrationRepository) MarkCompleted(ctx context.Context, token uuid.UUID) error {
	query := `UPDATE registration_requests SET status = 'completed' WHERE token = $1`
	_, err := r.pool.Exec(ctx, query, token)
	return err
}

func (r *RegistrationRepository) DeleteByToken(ctx context.Context, token uuid.UUID) error {
	query := `DELETE FROM registration_requests WHERE token = $1`
	_, err := r.pool.Exec(ctx, query, token)
	return err
}

func (r *RegistrationRepository) CleanExpired(ctx context.Context) error {
	query := `UPDATE registration_requests 
              SET status = 'expired' 
              WHERE status = 'pending' AND expires_at < NOW()`
	_, err := r.pool.Exec(ctx, query)
	return err
}
