package postgres

import (
	"context"
	"errors"
	"fmt"
	"study/internal/core/domain"

	"github.com/jackc/pgx/v5"
)

func (u *UserRepository) SaveUser(ctx context.Context, user *domain.User) (domain.User, error) {
	if user == nil {
		return domain.User{}, errors.New("user is nil")
	}

	query := `INSERT INTO users (login, password, role) VALUES ($1, $2, $3)
	RETURNING id`

	err := u.pool.QueryRow(ctx, query, user.Login, user.Password, user.Role).Scan(&user.ID)

	if err != nil {
		return domain.User{}, err
	}

	return *user, nil

}

func (u *UserRepository) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	if login == "" {
		return domain.User{}, fmt.Errorf("Login can`t be empity")
	}
	query := `SELECT id, login, password, role FROM users WHERE login = $1`
	var user domain.User

	err := u.pool.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with login %s not found", login)
		}
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
