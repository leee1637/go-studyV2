package service

import (
	"context"
	"errors"
	"fmt"
	"study/internal/core/domain"
	repository_postgres "study/internal/features/users/repository/postgres"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository_postgres.UserRepository
	secretKey []byte
}

func NewAuthService(repo *repository_postgres.UserRepository, secretKey string) *AuthService {
	return &AuthService{
		repo:      repo,
		secretKey: []byte(secretKey),
	}
}

func (s *AuthService) SignUp(ctx context.Context, login, password string, role domain.Role) error {

	_, err := s.repo.GetByLogin(ctx, login)

	if err == nil {
		return fmt.Errorf("Пользователь с логином %s уже существует", login)
	}

	if len(password) < 8 {
		return fmt.Errorf("Пароль не может быть меньше 8 символов")
	}

	if len(password) > 100 {
		return fmt.Errorf("Пароль не может быть больше 100 символов")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("Ошибка при создании bcrypt ключа: %W ", err)
	}

	NewUser := domain.User{
		Login:    login,
		Password: string(hashedPassword),
		Role:     role,
	}

	err = NewUser.Validate()
	if err != nil {
		return fmt.Errorf("Не прошёл валидацию: %w", err)
	}

	_, err = s.repo.SaveUser(ctx, &NewUser)

	if err != nil {
		return fmt.Errorf("Ошибка запроса бд: %w", err)
	}

	return nil

}

func (s *AuthService) SignIn(ctx context.Context, login, password string) (string, error) {
	// 1. Найти пользователя по логину через репозиторий
	user, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		return "", errors.New("Неверный логин")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", fmt.Errorf("Неверный пароль: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"login":   user.Login,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 часа
		"iat":     time.Now().Unix(),                     // время создания
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("Не удаётся сгенерировать токен: %w", err)
	}

	return tokenString, nil

}
