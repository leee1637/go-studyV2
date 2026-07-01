package students_postgres_repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (u *UserRepository) Begin(ctx context.Context) (pgx.Tx, error) {
	return u.pool.Begin(ctx) // u.pool — это твой *pgxpool.Pool внутри репозитория
}
