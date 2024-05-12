package request

import (
	"encoding/json"
	"net/http"
)

func DecodeJson(w http.ResponseWriter, r *http.Request, dest any) error {
	maxBytes := 1 << 20 // 1_048_576 - 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	return dec.Decode(dest)
}
