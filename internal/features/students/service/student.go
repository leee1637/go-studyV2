package student_service

import (
	"context"
	"fmt"
	"study/internal/core/domain"
	students_postgres_repository "study/internal/features/students/repository/student"
)

type StudentService struct {
	repo      *students_postgres_repository.UserRepository
	secretKey []byte
}

func NewAuthService(repo *students_postgres_repository.UserRepository, secretKey string) *AuthService {
	return &StudentService{
		repo:      repo,
		secretKey: []byte(secretKey),
	}
}

func (s *StudentService) GetAll(ctx context.Context, pag domain.PaginationRequest) (*PageResulte, error) {
	n := domain.NewPaginationRequest(pag.Page, pag.PageSize)

	request, err := s.repo.GetAll(ctx, n)

	if err != nil {
		return nil, fmt.Errorf("Ошибка получения запроса: %w", err)
	}

	return request, nil
}
