package student_service

import (
	"context"
	"fmt"
	"study/internal/core/domain"
	students_postgres_repository "study/internal/features/students/repository/postgres"
)

type StudentService struct {
	repo      *students_postgres_repository.StudentRepository
	secretKey []byte
}

func NewStudentService(repo *students_postgres_repository.StudentRepository, secretKey string) *StudentService {
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

func (s *StudentService) GetByIDStudent(ctx context.Context, idRequest int, role domain.Role, idOut int) (domain.Student, error) {
	switch role {
	case domain.RoleAdmin:
		s, err := s.repo.GetByID(ctx, idOut)
		if err != nil {
			return domain.Student{}, fmt.Errorf("Ошибка при вызова студентна: %w", err)
		}

		return s, err
	case domain.RoleTeacher:
		slGroup, err := s.repo.IsStudentInTeacherGroup(ctx, idRequest, idOut)
		if err != nil {
			return domain.Student{}, fmt.Errorf("Ошибка запроса поиска учителя и студента: %w", err)
		}

		if slGroup != true {
			return domain.Student{}, fmt.Errorf("Студент не состоит в группе запроса от учителя")
		}

		s, err := s.repo.GetByID(ctx, idOut)
		if err != nil {
			return domain.Student{}, fmt.Errorf("Ошибка при вызова студентна: %w", err)
		}

		return s, err

	case domain.RoleStudent:
		if idRequest != idOut {
			return domain.Student{}, fmt.Errorf("Нету достпуа, студент может смотреть тольк осебя")
		}

		s, err := s.repo.GetByID(ctx, idOut)
		if err != nil {
			return domain.Student{}, fmt.Errorf("Ошибка при вызова студентна: %w", err)
		}

		return s, err

	default:
		return domain.Student{}, fmt.Errorf("Ошибка сравнения ролей")
	}

}

func (s *StudentService) UpdateStudent(ctx context.Context, req int, role domain.Role, out int) error {
	ds, err := s.repo.GetByID(ctx, out)
	if err != nil {
		return fmt.Errorf("Ошибка при получении студентна: %w", err)
	}
	switch role {
	case domain.RoleAdmin:
		err := s.repo.UpdateStudent(ctx, ds)
		if err != nil {
			return fmt.Errorf("Ошибка при обновлении студентна: %w", err)
		}

		return nil
	case domain.RoleTeacher:
		slGroup, err := s.repo.IsTeacherOfGroup(ctx, req, ds.GroupName)
		if err != nil {
			return fmt.Errorf("Ошибка запроса поиска учителя и студента: %w", err)
		}

		if slGroup != true {
			return fmt.Errorf("Студент не состоит в группе запроса от учителя")
		}

		err = s.repo.UpdateStudent(ctx, ds)
		if err != nil {
			return fmt.Errorf("Ошибка при обновлении студентна: %w", err)
		}

		return nil

	case domain.RoleStudent:
		if req != ds.ID {
			return fmt.Errorf("Нету достпуа, студент может смотреть тольк осебя")
		}

		err := s.repo.UpdateStudent(ctx, ds)
		if err != nil {
			return fmt.Errorf("Ошибка при обновлении студентна: %w", err)
		}

		return nil

	default:
		return fmt.Errorf("Ошибка сравнения ролей")
	}

}

func (s *StudentService) GetByGroup(ctx context.Context, idReq int, role domain.Role, group string) ([]domain.Student, error) {
	switch role {
	case domain.RoleAdmin:
		s, err := s.repo.GetByGroup(ctx, group)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при вызова студентна по группе: %w", err)
		}

		return s, err
	case domain.RoleTeacher:
		slGroup, err := s.repo.IsTeacherOfGroup(ctx, idReq, group)
		if err != nil {
			return nil, fmt.Errorf("Ошибка запроса учителя к группе: %w", err)
		}
		if slGroup != true {
			return nil, fmt.Errorf("Студент не состоит в группе запроса от учителя")
		}

		s, err := s.repo.GetByGroup(ctx, group)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при вызова студентна по группе: %w", err)
		}

		return s, err

	case domain.RoleStudent:
		student, err := s.repo.GetByID(ctx, idReq)
		if err != nil {
			return nil, fmt.Errorf("Ошибка поиска по id: %w", err)
		}

		if student.GroupName != group {
			return nil, fmt.Errorf("Ошибка! Студент смотрит не свою группу!")
		}

		s, err := s.repo.GetByGroup(ctx, group)
		if err != nil {
			return nil, fmt.Errorf("Ошибка при вызова студентна по группе: %w", err)
		}

		return s, err

	default:
		return nil, fmt.Errorf("Ошибка сравнения ролей")
	}
}

func (s *StudentService) DeleteStudent(ctx context.Context, idReq int, role domain.Role, group string, idDel int) error {
	switch role {
	case domain.RoleAdmin:
		err := s.repo.DeleteStudent(ctx, idDel)
		if err != nil {
			return fmt.Errorf("Ошибка при удалении студентна: %w", err)
		}

		return nil
	case domain.RoleTeacher:
		slGroup, err := s.repo.IsStudentInTeacherGroup(ctx, idReq, idDel)
		if err != nil {
			fmt.Errorf("Ошибка запроса учителя к группе: %w", err)
		}
		if slGroup != true {
			fmt.Errorf("Студент не состоит в группе запроса от учителя")
		}

		err = s.repo.DeleteStudent(ctx, idDel)
		if err != nil {
			fmt.Errorf("Ошибка при удалении студентна: %w", err)
		}

		return nil

	case domain.RoleStudent:

		if idReq != idDel {
			return fmt.Errorf("Ошибка! Студент делитает не себя!")
		}

		err := s.repo.DeleteStudent(ctx, idDel)
		if err != nil {
			return fmt.Errorf("Ошибка при удалении студентна: %w", err)
		}

		return nil

	default:
		return fmt.Errorf("Ошибка сравнения ролей")
	}
}
