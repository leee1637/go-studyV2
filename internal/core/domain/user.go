package domain

import (
	"errors"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"-"`
	Role     Role   `json:"role"`
}

func NewUser(
	id int,
	login string,
	password string,
	role Role,
) *User {
	return &User{
		ID:       id,
		Login:    login,
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
