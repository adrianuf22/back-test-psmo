package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adrianuf22/back-test-psmo/internal"
	"github.com/adrianuf22/back-test-psmo/internal/domain/transaction"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	body := &transaction.Input{AccountID: int64(1), OperationTypeID: int(transaction.Withdrawal), Amount: 24.5}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	s := internal.NewServer(context.Background())
	s.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	dec := json.NewDecoder(rr.Body)

	actual := &struct {
		transaction.Output `json:"data"`
	}{transaction.Output{}}
	err = dec.Decode(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int64(1), actual.ID)
	assert.Equal(t, int64(1), actual.AccountID)
	assert.Equal(t, transaction.Withdrawal, actual.OperationTypeID)
	assert.Equal(t, -24.5, actual.Amount)
	assert.IsType(t, time.Now(), actual.EventDate)
}

func TestTransactionCustomerNotFound(t *testing.T) {
	body := &transaction.Input{AccountID: 1000, OperationTypeID: 3, Amount: 24.5}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	s := internal.NewServer(context.Background())
	s.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, `{"error":"not found"}`, rr.Body.String())
}

func TestTransactionInvalidOperationType(t *testing.T) {
	body := &transaction.Input{AccountID: 1, OperationTypeID: 7, Amount: 24.5}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	s := internal.NewServer(context.Background())
	s.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"error":"bad input data","values":{"operation_type_id":"invalid operation type id"}}`, rr.Body.String())
}
