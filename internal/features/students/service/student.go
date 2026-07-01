package student_service

import (
	"context"
	"fmt"
	"study/internal/core/domain"
	students_postgres_repository "study/internal/features/students/repository/postgres"
)

type StudentService struct {
	repo      *students_postgres_repository.UserRepository
	secretKey []byte
}

func NewStudentService(repo *students_postgres_repository.UserRepository, secretKey string) *StudentService {
	return &StudentService{
		repo:      repo,
		secretKey: []byte(secretKey),
	}
}

func (s *StudentService) GetAll(ctx context.Context, page, limit int) (*domain.PageResult, error) {
	n := domain.NewPaginationRequest(page, limit)

	request, err := s.repo.GetAll(ctx, n)

	if err != nil {
		return nil, fmt.Errorf("Ошибка получения запроса: %w", err)
	}

	return request, nil
}
