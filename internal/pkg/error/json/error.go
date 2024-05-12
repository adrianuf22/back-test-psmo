package json

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/adrianuf22/back-test-psmo/internal/pkg/error/api"
)

func WrapError(err error) error {
	shouldHandle := errors.Is(err, &json.UnmarshalTypeError{}) ||
		errors.Is(err, &json.SyntaxError{}) ||
		strings.Contains(err.Error(), "json:")

	if shouldHandle {
		return api.WrapError(err, api.ErrBadRequest)
	}

	return err
}
