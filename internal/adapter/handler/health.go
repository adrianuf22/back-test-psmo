package handler

import (
	"context"
	"net/http"

	"github.com/adrianuf22/back-test-psmo/internal/domain/health"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/request"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/response"
)

func RegisterHealthHandler(ctx context.Context, router *http.ServeMux, u health.Usecase) {
	v1 := "/v1/health"

	liveness := func(w http.ResponseWriter, r *http.Request) {
		response.Json(w, http.StatusOK, u.GetLivenessStatus())
	}
	router.HandleFunc(request.Get.WithPath(v1), liveness)
	router.HandleFunc(request.Get.WithPath(v1, "/liveness"), liveness)
	router.HandleFunc(request.Get.WithPath(v1, "/readiness"), func(w http.ResponseWriter, r *http.Request) {
		response.Json(w, http.StatusOK, u.GetReadinessStatus())
	})
}
