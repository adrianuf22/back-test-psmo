package postgres

import "database/sql"

type healthRepo struct {
	db *sql.DB
}

func NewHealthRepo(db *sql.DB) *healthRepo {
	return &healthRepo{
		db: db,
	}
}

func (r *healthRepo) Readiness() error {
	return r.db.Ping()
}
