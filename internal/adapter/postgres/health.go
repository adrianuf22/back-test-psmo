package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type healthRepo struct {
	db *pgxpool.Conn
}

func NewHealthRepo(db *pgxpool.Conn) *healthRepo {
	return &healthRepo{
		db: db,
	}
}

func (r *healthRepo) Readiness(ctx context.Context) error {
	return r.db.Ping(ctx)
}
