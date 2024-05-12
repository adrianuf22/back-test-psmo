package db

import (
	"database/sql"
	"errors"

	"github.com/adrianuf22/back-test-psmo/internal/pkg/error/api"
	"github.com/lib/pq"
)

func WrapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return api.WrapError(err, api.ErrNotFound)
	}

	pqErr := &pq.Error{}
	if !errors.As(err, &pqErr) {
		return err

	}
	switch {
	case pqErr.Code == "23505":
		return api.WrapError(err, api.ErrDuplicate)
	default:
		return err
	}
}
