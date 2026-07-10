package service_teacher

import (
	"context"
	"fmt"
	"study/internal/core/domain"
	repository_teacher "study/internal/features/teachers/repository"
)

type TeacherService struct {
	repo      *repository_teacher.TeacherRepository
	secretKey []byte
}

func NewTeacherService(repo *repository_teacher.TeacherRepository, secretKey string) *TeacherService {
	return &TeacherService{
		repo:      repo,
		secretKey: []byte(secretKey),
	}
}

func (s *TeacherService) GetAll(ctx context.Context, page, limit int) (*domain.PageResult, error) {
	n := domain.NewPaginationRequest(page, limit)

	request, err := s.repo.GetAll(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения запроса: %w", err)
	}

	return request, nil
}

func (s *TeacherService) GetByID(ctx context.Context, idReq, idResp int, role domain.Role) (domain.Teacher, error) {
	switch role {
	case domain.RoleAdmin:
		s, err := s.repo.GetByID(ctx, idResp)
		if err != nil {
			return domain.Teacher{}, fmt.Errorf("Ошибка запроса: %w", err)
		}

		return s, nil

	case domain.RoleTeacher:
		if idReq != idResp {
			return domain.Teacher{}, fmt.Errorf("Учитель может смотреть только себя")
		}

		s, err := s.repo.GetByID(ctx, idResp)
		if err != nil {
			return domain.Teacher{}, fmt.Errorf("Ошибка запроса: %w", err)
		}

		return s, nil
	default:
		return domain.Teacher{}, fmt.Errorf("Ошибка получения роли")
	}

}

func (s *TeacherService) UpdateTeacher(ctx context.Context, idReq, idResp int, role domain.Role, newTec domain.UpdateTeachertDTO) error {
	tec, err := s.repo.GetByID(ctx, idResp)
	if err != nil {
		return fmt.Errorf("Ошибка поиска пользователя: %w", err)
	}
	tec.FIO = newTec.FIO
	tec.PhoneNumber = *newTec.PhoneNumber

	switch role {
	case domain.RoleAdmin:
		err := s.repo.UpdateTeacher(ctx, tec)
		if err != nil {
			return fmt.Errorf("Ошибка запроса: %w", err)
		}

		return nil

	case domain.RoleTeacher:
		if idReq != idResp {
			return fmt.Errorf("Пользователя вызывает не сам себя!")
		}

		err := s.repo.UpdateTeacher(ctx, tec)
		if err != nil {
			return fmt.Errorf("Ошибка обновления учителя: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("Ошибка получения роли")
	}
}

func (s *TeacherService) DeleteTeacher(ctx context.Context, idReq, idResp int, role domain.Role) error {
	switch role {
	case domain.RoleAdmin:
		err := s.repo.DeleteTeacher(ctx, idResp)
		if err != nil {
			return fmt.Errorf("Ошибка запроса: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("Ошибка получения роли")
	}
}

func (s *TeacherService) GetByGroup(ctx context.Context, idReq int, group string, role domain.Role) ([]domain.Teacher, error) {
	switch role {
	case domain.RoleAdmin:
		tec, err := s.repo.GetByGroup(ctx, group)
		if err != nil {
			return nil, fmt.Errorf("Ошибка запроса: %w", err)
		}

		return tec, nil

	case domain.RoleTeacher:
		tec, err := s.repo.GetByID(ctx, idReq)
		if err != nil {
			return nil, fmt.Errorf("Ошибка поулчения препрда Id: %w", err)
		}
		for _, v := range tec.GroupName {
			if v == group {
				tecZ, err := s.repo.GetByGroup(ctx, group)
				if err != nil {
					return nil, fmt.Errorf("Ошибка запроса к группе: %w", err)
				}
				return tecZ, nil
			}
		}
		return nil, fmt.Errorf("Группа не преподователя")
	default:
		return nil, fmt.Errorf("Ошибка получения роли")
	}
}

func (s *TeacherService) AddGroup(ctx context.Context, idResp int, group string, role domain.Role) error {
	switch role {
	case domain.RoleAdmin:
		err := s.repo.AddGroup(ctx, idResp, group)
		if err != nil {
			return fmt.Errorf("Ошибка запроса: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("Ошибка получения роли")
	}
}

func (s *TeacherService) DeleteGroup(ctx context.Context, idResp int, group string, role domain.Role) error {
	switch role {
	case domain.RoleAdmin:
		err := s.repo.DeleteGroup(ctx, idResp, group)
		if err != nil {
			return fmt.Errorf("Ошибка запроса: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("Ошибка получения роли")
	}
}
