package students_postgres_repository

import (
	"context"
	"fmt"
	"study/internal/core/domain"
)

const (
	countQuery = `SELECT COUNT(*) FROM students`
)

func (u *StudentRepository) GetAll(ctx context.Context, pag domain.PaginationRequest) (*domain.PageResult, error) {
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

func (u *StudentRepository) GetByID(ctx context.Context, id int) (domain.Student, error) {
	if id < 0 {
		return domain.Student{}, fmt.Errorf("id не может быть меньше нуля")
	}
	query := `SELECT id, fio, group_name, phone_number FROM students
	WHERE id = $1`

	var s domain.Student

	err := u.pool.QueryRow(ctx, query, id).Scan(&s.ID, &s.FIO, &s.GroupName, &s.PhoneNumber)
	if err != nil {
		return domain.Student{}, fmt.Errorf("Ошибка запроста студента в бд", err)
	}

	return s, nil

}

func (u *StudentRepository) GetByGroup(ctx context.Context, group string) ([]domain.Student, error) {
	if group == "" {
		return nil, fmt.Errorf("группа не может быть пустая")
	}
	query := `SELECT id, fio, group_name, phone_number FROM students
	WHERE group_name = $1`

	rows, err := u.pool.Query(ctx, query, group)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса в бд: %w", err)
	}

	var students []domain.Student

	for rows.Next() {
		var s domain.Student
		err := rows.Scan(&s.ID, &s.FIO, &s.GroupName, &s.PhoneNumber)
		if err != nil {
			return nil, fmt.Errorf("Ошибка записи данных с запроса: %w", err)
		}

		students = append(students, s)
	}

	return students, nil
}

func (u *StudentRepository) UpdateStudent(ctx context.Context, s domain.Student) error {
	query := ` UPDATE students SET fio=$2, group_name=$3, phone_number=$4 WHERE id=$1`

	row, err := u.pool.Exec(ctx, query, s.ID, s.FIO, s.GroupName, s.PhoneNumber)
	if err != nil {
		return fmt.Errorf("Ошибка обновления пользвоателя: %w", err)
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("Пользователь не найден: %w", err)
	}

	return nil
}

func (u *StudentRepository) DeleteStudent(ctx context.Context, id int) error {
	if id < 0 {
		return fmt.Errorf("id не может быть меньше нуля")
	}
	query := `DELETE FROM students WHERE id = $1`

	row, err := u.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Ошибка запроса: %w", err)
	}

	if row.RowsAffected() == 0 {
		return fmt.Errorf("Юзер не найден: %w", err)
	}

	return nil
}

func (u *StudentRepository) IsStudentInTeacherGroup(ctx context.Context, teacherID, userID int) (bool, error) {
	query := `SELECT EXISTS (
	SELECT 1 FROM teachers_group tg
	JOIN students s ON s.group_name = tg.group_name
	WHERE tg.teacher_id = $1 AND s.id = $2)`

	var exists bool
	err := u.pool.QueryRow(ctx, query, teacherID, userID).Scan(&exists)
	return exists, err
}

func (u *StudentRepository) IsTeacherOfGroup(ctx context.Context, teacherID int, gr string) (bool, error) {
	query := `SELECT EXISTS (
	SELECT 1 FROM teachers_group tg
	WHERE teacher_id = $1 AND group_name = $2)`

	var exists bool
	err := u.pool.QueryRow(ctx, query, teacherID, gr).Scan(&exists)
	return exists, err
}
