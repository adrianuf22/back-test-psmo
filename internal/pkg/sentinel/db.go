package sentinel

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

func WrapDBError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return WrapError(err, ErrNotFound)
	}

	pqErr := &pq.Error{}
	if !errors.As(err, &pqErr) {
		return err

	}
	switch {
	case pqErr.Code == "23505":
		return WrapError(err, ErrDuplicate)
	default:
		return err
	}
}
