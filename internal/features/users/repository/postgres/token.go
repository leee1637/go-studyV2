package repository_postgres

import (
	"context"
	"fmt"
	"study/internal/core/domain"
	"time"

	"github.com/jackc/pgx/v5"
)

func (u *UserRepository) SaveRefreshTokenNoTX(ctx context.Context,
	userID int,
	token string,
	expiresAt time.Time,
) error {
	query := `INSERT INTO refresh_tokens (id_users, token, expires_at, created_at) 
	VALUES ($1, $2, $3, NOW())`
	_, err := u.pool.Exec(ctx, query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("Ошибка создания токена: %w", err)
	}

	return nil
}

func (u *UserRepository) SaveRefreshToken(ctx context.Context, tx pgx.Tx,
	userID int,
	token string,
	expiresAt time.Time,
) error {
	query := `INSERT INTO refresh_tokens (id_users, token, expires_at, created_at) 
	VALUES ($1, $2, $3, NOW())`
	_, err := tx.Exec(ctx, query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("Ошибка создания токена: %w", err)
	}

	return nil
}

func (u *UserRepository) GetRefreshToken(ctx context.Context, tx pgx.Tx, token string) (domain.RefreshToken, error) {
	query := `SELECT id, id_users, token, expires_at, created_at
	FROM refresh_tokens WHERE token=$1`

	var rt domain.RefreshToken

	err := u.pool.QueryRow(ctx, query, token).Scan(&rt.ID, &rt.IDUser, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return domain.RefreshToken{}, fmt.Errorf("Ошибка получения токена: %w", err)
	}

	return rt, nil
}
func (u *UserRepository) DeleteRefreshTokenNoTX(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token=$1`

	_, err := u.pool.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("Ошибка удаления пользователя по токену: %w", err)
	}

	return nil
}

func (u *UserRepository) DeleteRefreshToken(ctx context.Context, tx pgx.Tx, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token=$1`

	_, err := tx.Exec(ctx, query, token)
	if err != nil {
		return fmt.Errorf("Ошибка удаления пользователя по токену: %w", err)
	}

	return nil
}

func (u *UserRepository) DeleteAllRefreshToken(ctx context.Context, idUser int) error {
	query := `DELETE FROM refresh_tokens WHERE id_users=$1`

	_, err := u.pool.Exec(ctx, query, idUser)
	if err != nil {
		return fmt.Errorf("Ошибка удаления пользователя по айди: %w", err)
	}

	return nil
}
