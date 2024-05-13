package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/adrianuf22/back-test-psmo/internal/domain/transaction"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/sentinel"
	"github.com/stretchr/testify/assert"
)

type MockedTransactionService struct {
	dummy *transaction.Model
	err   error
}

func (m *MockedTransactionService) Create(model *transaction.Model) error {
	model.SetID(m.dummy.ID())

	return m.err
}

func TestCreateTransaction(t *testing.T) {
	accountID, operationTypeID, amount := int64(1), 3, 24.5

	body := &transaction.Input{AccountID: accountID, OperationTypeID: operationTypeID, Amount: amount}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	dummyModel := &transaction.Model{}
	dummyModel.SetID(10)

	dummyAccountModel := &account.Model{}
	dummyAccountModel.SetID(1)

	handler := transactHandler{
		ctx:     context.Background(),
		usecase: transaction.NewUsecase(&MockedTransactionService{dummy: dummyModel}, &MockedAccountService{dummy: dummyAccountModel}),
	}

	http.HandlerFunc(handler.Create).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	dec := json.NewDecoder(rr.Body)

	actual := &struct {
		transaction.Output `json:"data"`
	}{transaction.Output{}}
	err = dec.Decode(actual)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, int64(10), actual.ID)
	assert.Equal(t, accountID, actual.AccountID)
	assert.Equal(t, transaction.Withdrawal, actual.OperationTypeID)
	assert.Equal(t, amount, actual.Amount)
	assert.IsType(t, time.Now(), actual.EventDate)
}

func TestTransactionCustomerNotFound(t *testing.T) {
	body := &transaction.Input{AccountID: 1, OperationTypeID: 3, Amount: 24.5}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/transactions", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	dummyModel := &transaction.Model{}
	dummyModel.SetID(10)

	handler := transactHandler{
		ctx:     context.Background(),
		usecase: transaction.NewUsecase(&MockedTransactionService{dummy: dummyModel}, &MockedAccountService{dummy: nil, err: sentinel.ErrNotFound}),
	}

	http.HandlerFunc(handler.Create).ServeHTTP(rr, req)

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

	dummyModel := &transaction.Model{}
	dummyModel.SetID(10)

	handler := transactHandler{
		ctx:     context.Background(),
		usecase: transaction.NewUsecase(&MockedTransactionService{dummy: dummyModel}, &MockedAccountService{dummy: nil, err: sentinel.ErrNotFound}),
	}

	http.HandlerFunc(handler.Create).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, `{"error":"bad input data","values":{"operation_type_id":"invalid operation type id"}}`, rr.Body.String())
}
