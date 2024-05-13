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

type accountHandler struct {
	ctx     context.Context
	usecase account.Usecase
}

func RegisterAccountHandler(ctx context.Context, router *http.ServeMux, u account.Usecase) {
	h := &accountHandler{
		ctx:     ctx,
		usecase: u,
	}

	v1 := "/v1/accounts"
	router.HandleFunc(request.Get.WithPath(v1, "/{accountId}"), h.Read)
	router.HandleFunc(request.Post.WithPath(v1), h.Create)
}

func (h *accountHandler) Read(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("accountId"))
	if err != nil {
		response.ErrorJson(w, sentinel.ErrBadRequest)
		log.Println(err)
		return
	}

	// TODO Adjust timeout using http config
	ctx, cancel := context.WithTimeout(h.ctx, 3*time.Second)
	defer cancel()

	found, err := h.usecase.GetAccountById(ctx, int64(id))
	if err != nil {
		response.ErrorJson(w, err)
		return
	}

	response.Json(w, http.StatusOK, account.ToOutput(found))
}

func (h *accountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input account.Input
	err := request.DecodeJson(w, r, &input)
	if err != nil {
		response.ErrorJson(w, err)
		return
	}

	created, err := h.usecase.CreateAccount(h.ctx, input)
	if err != nil {
		response.ErrorJson(w, err)
		return
	}

	response.Json(w, http.StatusCreated, account.ToOutput(created))
}
