package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adrianuf22/back-test-psmo/internal/domain/account"

	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/request"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/handler/response"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
)

type handler struct {
	ctx     context.Context
	usecase account.Usecase
}

func RegisterAccountHandler(ctx context.Context, router *http.ServeMux, u account.Usecase) {
	h := &handler{
		ctx:     ctx,
		usecase: u,
	}

	v1 := "/v1/account"
	router.HandleFunc(request.Get.WithPath(v1, "/{accountId}"), h.Account)
	router.HandleFunc(request.Post.WithPath(v1), h.Create)
}

func (h *handler) Account(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("accountId"))
	if err != nil {
		response.ErrorJson(w, sentinel.ErrBadRequest)
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
