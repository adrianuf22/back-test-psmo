package handler

import (
	"context"
	"net/http"

	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/response"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
)

func RegisterErrorHandler(ctx context.Context, router *http.ServeMux) {
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		response.ErrorJson(w, sentinel.ErrNotFound)
	})
}
