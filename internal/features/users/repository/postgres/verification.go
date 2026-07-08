package repository_postgres

import (
	"context"
	"errors"
	"fmt"
	"study/internal/core/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (u *UserRepository) SaveVerificationToken(ctx context.Context, tx pgx.Tx, v *domain.EmailVerification) error {
	query := `INSERT INTO email_verification (user_id, token, expires_at)
	VALUES ($1, $2, $3)`

	_, err := tx.Exec(ctx, query, v.UserID, v.Token, v.ExpiresAt)
	if err != nil {
		return fmt.Errorf("Ошибка запроса к бд: %w", err)
	}

	return nil
}

func (u *UserRepository) SaveVerificationTokenNoTX(ctx context.Context, v domain.EmailVerification) error {
	query := `INSERT INTO email_verification (user_id, token, expires_at)
	VALUES ($1, $2, $3)`

	_, err := u.pool.Exec(ctx, query, v.UserID, v.Token, v.ExpiresAt)
	if err != nil {
		return fmt.Errorf("Ошибка запроса к бд: %w", err)
	}

	return nil
}

func (u *UserRepository) GetVerificationByToken(ctx context.Context, token uuid.UUID) (domain.EmailVerification, error) {
	if len(token) == 0 {
		return domain.EmailVerification{}, fmt.Errorf("Пустой токен")
	}

	query := `SELECT id, user_id, token, expires_at, created_at 
	FROM email_verification WHERE token=$1`

	var de domain.EmailVerification

	err := u.pool.QueryRow(ctx, query, token).Scan(&de.ID, &de.UserID, &de.Token, &de.ExpiresAt, &de.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.EmailVerification{}, fmt.Errorf("Токен не найден")
		}
		return domain.EmailVerification{}, fmt.Errorf("Ошибка БД: %w", err)
	}

	return de, nil
}

func (u *UserRepository) MarkEmailVerified(ctx context.Context, tx pgx.Tx, userID int) error {
	query := `UPDATE users SET is_verified = TRUE WHERE id=$1`

	_, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("Ошибка обновлени пользовтаеля: %w", err)
	}

	return nil
}

func (u *UserRepository) DeleteVerificationToken(ctx context.Context, tx pgx.Tx, token uuid.UUID) error {
	query := `DELETE FROM email_verification WHERE token=$1`

	row, err := tx.Exec(ctx, query, token)

	if err != nil {
		return fmt.Errorf("ошибка БД при удалении токена: %w", err)
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("токен не найден")
	}

	return nil
}

func (u *UserRepository) IsEmailVerified(ctx context.Context, userID int) (bool, error) {
	query := `SELECT is_verified FROM users WHERE id=$1`

	var ch bool

	err := u.pool.QueryRow(ctx, query, userID).Scan(&ch)

	return ch, err
}
