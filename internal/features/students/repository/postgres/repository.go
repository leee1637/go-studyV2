package students_postgres_repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{
		pool: pool,
	}
}

func (u *StudentRepository) Begin(ctx context.Context) (pgx.Tx, error) {
	return u.pool.Begin(ctx) // u.pool — это твой *pgxpool.Pool внутри репозитория
}
