package repository_teacher

import "github.com/jackc/pgx/v5/pgxpool"

type TeacherRepository struct {
	pool *pgxpool.Pool
}

func NewTeacherRepository(pool *pgxpool.Pool) *TeacherRepository {
	return &TeacherRepository{
		pool: pool,
	}
}
