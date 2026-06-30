package repository_postgres

import (
	"context"
	"errors"
	"fmt"
	"study/internal/core/domain"

	"github.com/jackc/pgx/v5"
)

func (u *UserRepository) SaveUser(ctx context.Context, tx pgx.Tx, user *domain.User) (int, error) {
	if user == nil {
		return 0, errors.New("user is nil")
	}

	query := `INSERT INTO users (login, password, role) VALUES ($1, $2, $3)
	RETURNING id`

	err := tx.QueryRow(ctx, query, user.Login, user.Password, user.Role).Scan(&user.ID)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (u *UserRepository) SaveUserStudent(ctx context.Context, tx pgx.Tx, user *domain.Student) error {
	if user == nil {
		return errors.New("user is nil")
	}

	query := `INSERT INTO students (id, fio, group_name, phone_number) VALUES ($1, $2, $3, $4)`

	_, err := tx.Exec(ctx, query, user.ID, user.FIO, user.GroupName, user.PhoneNumber)

	return err

}

func (u *UserRepository) SaveUserTeacher(ctx context.Context, tx pgx.Tx, user *domain.Teacher) error {
	if user == nil {
		return errors.New("user is nil")
	}

	query := `INSERT INTO teachers (id, fio, phone_number) VALUES ($1, $2, $3)`

	_, err := tx.Exec(ctx, query, user.ID, user.FIO, user.PhoneNumber)
	if err != nil {
		return errors.New("Ошибка добавления препода запроса")
	}

	if len(user.GroupName) > 0 {
		query := `INSERT INTO teachers_group (teacher_id, group_name) VALUES ($1, $2)`

		for _, v := range user.GroupName {
			_, err := tx.Exec(ctx, query, user.ID, v)
			if err != nil {
				return errors.New("Ошибка привязки к группе")
			}
		}
	}
	return nil
}

func (u *UserRepository) SaveUserAdmin(ctx context.Context, tx pgx.Tx, user *domain.Admin) error {
	if user == nil {
		return errors.New("user is nil")
	}

	query := `INSERT INTO admins (id, fio, phon_number) VALUES ($1, $2, $3)`

	_, err := tx.Exec(ctx, query, user.ID, user.FIO, user.PhoneNumber)

	return err
}

func (u *UserRepository) GetByLogin(ctx context.Context, login string) (domain.SignUpDTO, error) {
	if login == "" {
		return domain.SignUpDTO{}, fmt.Errorf("Login can`t be empity")
	}
	query := `SELECT id, login, password, role FROM users WHERE login = $1`
	var user domain.SignUpDTO

	err := u.pool.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password, &user.Role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SignUpDTO{}, fmt.Errorf("user with login %s not found. err: %w", login, pgx.ErrNoRows)
		}
		return domain.SignUpDTO{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
