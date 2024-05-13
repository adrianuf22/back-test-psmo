package sentinel

import (
	"encoding/json"
	"errors"
	"strings"
)

func WrapJsonError(err error) error {
	shouldHandle := errors.Is(err, &json.UnmarshalTypeError{}) ||
		errors.Is(err, &json.SyntaxError{}) ||
		strings.Contains(err.Error(), "json:")

	if shouldHandle {
		return WrapError(err, ErrBadRequest)
	}

	return err
}
