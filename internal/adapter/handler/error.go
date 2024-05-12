package handler

import (
	"context"
	"net/http"

	"github.com/adrianuf22/back-test-psmo/internal/pkg/error/api"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/response"
	"github.com/go-chi/chi/v5"
)

func RegisterErrorHandler(ctx context.Context, router *chi.Mux) {
	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		response.ErrorJson(w, api.ErrNotFound)
	})
}
