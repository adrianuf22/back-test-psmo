package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
)

type body struct {
	Data interface{} `json:"data"`
}

func Json(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if payload == nil {
		payload = []byte("[]")
	}

	write(w, &body{
		Data: payload,
	}, statusCode)
}

func ErrorJson(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/problem+json")

	var httpErr sentinel.Error
	if !errors.As(err, &httpErr) {
		httpErr = sentinel.ErrInternal
	}

	status, msg := httpErr.HttpError()
	write(w, msg, status)
}

func write(w http.ResponseWriter, payload interface{}, statusCode int) {
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error(err.Error())
		writeFatalError(w)
		return
	}

	if string(data) == "null" {
		data = []byte("[]")
	}

	if _, err = w.Write(data); err != nil {
		slog.Error(err.Error())
		writeFatalError(w)
	}
}

func writeFatalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(sentinel.ErrSafeFatal)
}
