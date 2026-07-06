package repository_teacher

import (
	"context"
	"fmt"
	"study/internal/core/domain"
)

func (t *TeacherRepository) GetAll(ctx context.Context, pag domain.PaginationRequest) (*domain.PageResult, error) {
	query := `SELECT t.id, t.fio, t.phone_number, 
	ARRAY_AGG(tg.group_name)
	FROM teachers t
	LEFT JOIN teachers_group tg ON t.id = tg.teacher_id
	GROUP BY t.id, t.fio, t.phone_number
	ORDER BY t.id
	LIMIT $1 OFFSET $2`

	row, err := t.pool.Query(ctx, query, pag.GetLimit(), pag.GetOffset())
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса к бд: %w", err)
	}
	defer row.Close()

	Cquery := `SELECT COUNT(*) FROM teachers`

	var total int

	err = t.pool.QueryRow(ctx, Cquery).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса к бд на кол-во: %w", err)
	}

	var teachers []domain.Teacher

	for row.Next() {
		var teach domain.Teacher
		err := row.Scan(
			&teach.ID,
			&teach.FIO,
			&teach.PhoneNumber,
			&teach.GroupName,
		)
		if err != nil {
			return nil, fmt.Errorf("Ошибка импорта данных с бд: %w", err)
		}

		teachers = append(teachers, teach)
	}
	dn := domain.NewPageResult(teachers, pag, total)

	return &dn, nil

}

func (t *TeacherRepository) GetByID(ctx context.Context, id int) (domain.Teacher, error) {
	if id < 0 {
		return domain.Teacher{}, fmt.Errorf("id < 0 Ошибка")
	}
	query := `SELECT t.id, t.fio, t.phone_number, 
	ARRAY_AGG(tg.group_name)
	FROM teachers t
	LEFT JOIN teachers_group tg ON t.id = tg.teacher_id
	WHERE t.id = $1
	GROUP BY t.id, t.fio, t.phone_number`

	var tec domain.Teacher

	err := t.pool.QueryRow(ctx, query, id).Scan(&tec.ID, &tec.FIO, &tec.PhoneNumber, &tec.GroupName)
	if err != nil {
		return domain.Teacher{}, fmt.Errorf("ОШибка запроса к бд: %w", err)
	}

	return tec, nil
}

func (t *TeacherRepository) UpdateTeacher(ctx context.Context, tec domain.Teacher) error {
	query := `UPDATE teachers SET fio=$1, phone_number=$2
	WHERE id=$3`

	row, err := t.pool.Exec(ctx, query, tec.FIO, tec.PhoneNumber, tec.ID)
	if err != nil {
		return fmt.Errorf("Ошибка запроса: %w", err)
	}
	if row.RowsAffected() == 0 {
		return fmt.Errorf("Не найден учитель")
	}

	return nil
}

func (t *TeacherRepository) DeleteTeacher(ctx context.Context, id int) error {
	if id < 0 {
		return fmt.Errorf("id < 0 Ошибка")
	}
	query := `DELETE FROM teachers WHERE id=$1`

	row, err := t.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Ошибка запроса: %w", err)
	}
	if row.RowsAffected() == 0 {
		return fmt.Errorf("Не найден учитель")
	}

	return nil
}

func (t *TeacherRepository) GetByGroup(ctx context.Context, group string) ([]domain.Teacher, error) {
	if group == "" {
		return nil, fmt.Errorf("ошибка - группа пустая")
	}
	query := `SELECT t.id, t.fio, t.phone_number
	FROM teachers t
	INNER JOIN teachers_group tg ON t.id = tg.teacher_id
	WHERE tg.group_name = $1`

	var tec []domain.Teacher
	rows, err := t.pool.Query(ctx, query, group)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса: %w", err)
	}

	for rows.Next() {
		var tt domain.Teacher
		err = rows.Scan(&tt.ID, &tt.FIO, &tt.PhoneNumber)
		if err != nil {
			return nil, fmt.Errorf("Ошибка скана: %w", err)
		}
		tec = append(tec, tt)
	}

	return tec, nil
}

func (t *TeacherRepository) AddGroup(ctx context.Context, id int, group string) error {
	if group == "" {
		return fmt.Errorf("ошибка - группа пустая")
	}

	checkTeacherQuery := `SELECT id FROM teachers WHERE id = $1`
	row, err := t.pool.Exec(ctx, checkTeacherQuery, id)
	if err != nil {
		return fmt.Errorf("Ошибка запроса: %w", err)
	}
	if row.RowsAffected() == 0 {
		return fmt.Errorf("Не найден учитель")
	}

	query := `INSERT INTO teachers_group (group_name, teacher_id) 
	VALUES ($1, $2)`

	row, err = t.pool.Exec(ctx, query, group, id)
	if err != nil {
		return fmt.Errorf("Ошибка запроса: %w", err)
	}
	if row.RowsAffected() == 0 {
		return fmt.Errorf("Не удалось добавить группу")
	}

	return nil
}

func (t *TeacherRepository) DeleteGroup(ctx context.Context, group string) error {
	if group == "" {
		return fmt.Errorf("ошибка - группа пустая")
	}
	query := `DELETE FROM teachers_group WHERE group_name = $1`

	row, err := t.pool.Exec(ctx, query, group)
	if err != nil {
		return fmt.Errorf("Ошибка запроса: %w", err)
	}
	if row.RowsAffected() == 0 {
		return fmt.Errorf("Не найден учитель")
	}

	return nil
}
