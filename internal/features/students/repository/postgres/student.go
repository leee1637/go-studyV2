package students_postgres_repository

import (
	"context"
	"fmt"
	"study/internal/core/domain"
)

const (
	countQuery = `SELECT COUNT(*) FROM students`
)

func (u *UserRepository) GetAll(ctx context.Context, pag domain.PaginationRequest) (*domain.PageResult, error) {
	query := `SELECT id, fio, group_name, phone_number FROM students
	ORDER BY id
	LIMIT $1 OFFSET $2`

	rows, err := u.pool.Query(ctx, query, pag.GetLimit(), pag.GetOffset())
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса студента: %w", err)
	}

	defer rows.Close()

	var student []domain.Student

	for rows.Next() {
		var stud domain.Student

		err := rows.Scan(
			&stud.ID,
			&stud.FIO,
			&stud.GroupName,
			&stud.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("Ошибка записи студента: %w", err)
		}

		student = append(student, stud)
	}

	var total int

	err = u.pool.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при общем подсчёте: %w", err)
	}

	result := domain.NewPageResult(student, pag, total)
	return &result, nil

}
