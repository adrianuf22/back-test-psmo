package handler

import (
	"context"
	"net/http"

	"github.com/adrianuf22/back-test-psmo/internal/domain/transaction"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/request"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/response"
)

type transactHandler struct {
	ctx     context.Context
	usecase transaction.Usecase
}

func RegisterTransactionHandler(ctx context.Context, router *http.ServeMux, u transaction.Usecase) {
	h := &transactHandler{
		ctx:     ctx,
		usecase: u,
	}

	v1 := "/v1/transactions"
	router.HandleFunc(request.Post.WithPath(v1), h.Create)
}

func (h *transactHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input transaction.Input
	err := request.DecodeJson(w, r, &input)
	if err != nil {
		response.ErrorJson(w, err)
		return
	}

	created, err := h.usecase.CreateTransaction(h.ctx, input)

	if err != nil {
		response.ErrorJson(w, err)
		return
	}

	response.Json(w, http.StatusCreated, transaction.ToOutput(created))
}
