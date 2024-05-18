package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type healthRepo struct {
	db *pgx.Conn
}

func NewHealthRepo(db *pgx.Conn) *healthRepo {
	return &healthRepo{
		db: db,
	}
}

func (r *healthRepo) Readiness(ctx context.Context) error {
	return r.db.Ping(ctx)
}
