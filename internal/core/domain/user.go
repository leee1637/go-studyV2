package domain

import (
	"errors"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     Role   `json:"role"`
}

func NewUser(
	id int,
	email string,
	password string,
	role Role,
) *User {
	return &User{
		ID:       id,
		Email:    email,
		Password: password,
		Role:     role,
	}
}

func (u *User) Validate() error {
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("Пароль не может быть пустым")
	}

	if u.Role != RoleAdmin && u.Role != RoleStudent && u.Role != RoleTeacher {
		return errors.New("Неверна задана роль! Только STUDENT, ADMIN, TEACHER!")
	}

	return nil
}
