package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adrianuf22/back-test-psmo/internal"
	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	body := &account.Input{DocumentNumber: "1234567805"}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/accounts", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	s := internal.NewServer(context.Background())
	s.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, `{"data":{"account_id":2,"document_number":"1234567805"}}`, rr.Body.String())
}

func TestReadAccount(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/accounts/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	s := internal.NewServer(context.Background())
	s.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"data":{"account_id":1,"document_number":"9876543210"}}`, rr.Body.String())
}

func TestAccountNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/accounts/12345", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	s := internal.NewServer(context.Background())
	s.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, `{"error":"not found"}`, rr.Body.String())
}
