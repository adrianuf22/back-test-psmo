package request

import (
	"encoding/json"
	"net/http"

	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/validation"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
)

func DecodeJson(w http.ResponseWriter, r *http.Request, input validation.Input) error {
	maxBytes := 1 << 20 // 1_048_576 - 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	// dec.DisallowUnknownFields()

	err := dec.Decode(input)
	if err != nil {
		return sentinel.WrapJsonError(err)
	}

	errs := input.Validate()
	if len(errs) > 0 {
		return sentinel.ErrBadRequest.WithValues(errs)
	}

	return nil
}
