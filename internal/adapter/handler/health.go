package handler

import (
	"context"
	"net/http"

	"github.com/adrianuf22/back-test-psmo/internal/domain/health"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/response"
	"github.com/go-chi/chi/v5"
)

func RegisterHealthHandler(ctx context.Context, router *chi.Mux, u health.Usecase) {
	router.Route("/v1/health", func(r chi.Router) {
		liveness := func(w http.ResponseWriter, r *http.Request) {
			response.Json(w, http.StatusOK, u.GetLivenessStatus())
		}

		r.Get("/", liveness)
		r.Get("/liveness", liveness)
		r.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
			response.Json(w, http.StatusOK, u.GetReadinessStatus())
		})
	})
}
