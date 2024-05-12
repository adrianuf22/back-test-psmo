package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/error/api"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/error/json"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/request"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/response"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	ctx     context.Context
	usecase account.Usecase
}

func RegisterAccountHandler(ctx context.Context, router *chi.Mux, u account.Usecase) {
	h := &handler{
		ctx:     ctx,
		usecase: u,
	}

	router.Route("/v1/accounts", func(router chi.Router) {
		router.Get("/{accountId}", h.Account)
		router.Post("/", h.Create)
	})
}

func (h *handler) Account(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "accountId"))
	if err != nil {
		response.ErrorJson(w, api.ErrBadRequest)
		log.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(h.ctx, 3*time.Second)
	defer cancel()

	account, err := h.usecase.GetAccountById(ctx, id)
	if err != nil {
		response.ErrorJson(w, err)
		return
	}

	response.Json(w, http.StatusOK, account)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input account.Model
	err := request.DecodeJson(w, r, &input)
	if err != nil {
		response.ErrorJson(w, json.WrapError(err))
		return
	}

	errs := input.Validate()
	if len(errs) > 0 {
		response.ErrorJson(w, api.ErrBadRequest.WithValues(errs))
		return
	}

	account, err := h.usecase.CreateAccount(h.ctx, input)
	if err != nil {
		response.ErrorJson(w, err)
		return
	}

	response.Json(w, http.StatusCreated, account)
}
