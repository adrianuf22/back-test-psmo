package mux

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
	"github.com/stretchr/testify/assert"
)

type MockedAccountService struct {
	dummy *account.Model
	err   error
}

func (m *MockedAccountService) Read(ctx context.Context, id int64) (*account.Model, error) {
	return m.dummy, m.err
}

func (m *MockedAccountService) Create(ctx context.Context, model *account.Model) error {
	model.SetID(m.dummy.ID())

	return m.err
}

func TestCreateAccount(t *testing.T) {
	body := &account.Input{DocumentNumber: "1234567809"}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/accounts", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	dummyModel := &account.Model{}
	dummyModel.SetID(1)

	handler := accountHandler{
		ctx:     context.Background(),
		usecase: account.NewUsecase(&MockedAccountService{dummy: dummyModel}),
	}

	http.HandlerFunc(handler.Create).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, `{"data":{"account_id":1,"document_number":"1234567809"}}`, rr.Body.String())
}

func TestReadAccount(t *testing.T) {
	dummyModel := account.NewModel(1, "12345")

	handler := accountHandler{
		ctx:     context.Background(),
		usecase: account.NewUsecase(&MockedAccountService{dummy: dummyModel}),
	}

	req, err := http.NewRequest("GET", "/v1/accounts/{accountId}", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("accountId", "12345")

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Read)
	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"data":{"account_id":1,"document_number":"12345"}}`, rr.Body.String())
}

func TestAccountNotFound(t *testing.T) {
	handler := accountHandler{
		ctx:     context.Background(),
		usecase: account.NewUsecase(&MockedAccountService{dummy: nil, err: sentinel.ErrNotFound}),
	}

	req, err := http.NewRequest("GET", "/v1/accounts/{accountId}", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetPathValue("accountId", "12345")

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler.Read)
	h.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, `{"error":"not found"}`, rr.Body.String())
}
