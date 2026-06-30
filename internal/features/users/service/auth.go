package service

import (
	"context"
	"errors"
	"fmt"
	"study/internal/core/domain"
	repository_postgres "study/internal/features/users/repository/postgres"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"

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

type SignUpDTO struct {
	ID          int
	Login       string
	Password    string
	Role        domain.Role
	FIO         string
	GroupName   []string
	PhoneNumber *string
}

func (s *AuthService) SignUp(ctx context.Context, user SignUpDTO) error {

	_, err := s.repo.GetByLogin(ctx, user.Login)

	if err == nil {
		return fmt.Errorf("Такой логин уже есть!")
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("ошибка проверки логина: %w", err)
	}

	if len(user.Password) < 8 {
		return fmt.Errorf("Пароль не может быть меньше 8 символов")
	}

	if len(user.Password) > 100 {
		return fmt.Errorf("Пароль не может быть больше 100 символов")
	}

	if user.FIO == "" {
		return fmt.Errorf("Ошибка! Не указано ФИО пользователя!")
	}

	if user.Role == domain.RoleStudent && len(user.GroupName) == 0 {
		return fmt.Errorf("для роли СТУДЕНТ обязательно указание группы")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("Ошибка при создании bcrypt ключа: %W ", err)
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	NewUser := domain.User{
		Login:    user.Login,
		Password: string(hashedPassword),
		Role:     user.Role,
	}

	err = NewUser.Validate()
	if err != nil {
		return fmt.Errorf("Не прошёл валидацию: %w", err)
	}

	userID, err := s.repo.SaveUser(ctx, tx, &NewUser)
	if err != nil {
		return fmt.Errorf("ошибка сохранения базового пользователя в БД: %w", err)
	}

	switch user.Role {

	case domain.RoleStudent:

		newStudent := domain.Student{
			ID:          userID,
			FIO:         user.FIO,
			GroupName:   user.GroupName[0],
			PhoneNumber: user.PhoneNumber,
		}

		err = s.repo.SaveUserStudent(ctx, tx, &newStudent)
		if err != nil {
			return fmt.Errorf("ошибка сохранения профиля студента: %w", err)
		}

	case domain.RoleTeacher:

		newTeacher := domain.Teacher{
			ID:          userID,
			FIO:         user.FIO,
			PhoneNumber: *user.PhoneNumber,
			GroupName:   user.GroupName,
		}

		err = s.repo.SaveUserTeacher(ctx, tx, &newTeacher)
		if err != nil {
			return fmt.Errorf("ошибка сохранения профиля преподавателя: %w", err)
		}

	case domain.RoleAdmin:

		newAdmin := domain.Admin{
			ID:          userID,
			FIO:         user.FIO,
			PhoneNumber: user.PhoneNumber,
		}

		err = s.repo.SaveUserAdmin(ctx, tx, &newAdmin)
		if err != nil {
			return fmt.Errorf("ошибка сохранения профиля администратора: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
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
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("Не удаётся сгенерировать токен: %w", err)
	}

	return tokenString, nil

}
