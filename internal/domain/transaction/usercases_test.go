package transaction

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/adrianuf22/back-test-psmo/internal/domain/account"
	"github.com/adrianuf22/back-test-psmo/internal/pkg/atomic"
	"github.com/stretchr/testify/assert"
)

var (
	accountID    = int64(1000)
	accountDummy = &account.Model{}
)

func init() {
	accountDummy.SetID(accountID)
}

type MockedTransactionService struct {
	createMock           *Model
	readAllPurchasesMock []Model
	errMock              error
	callstack            []string
}

func (m *MockedTransactionService) Create(ctx context.Context, model *Model) error {
	m.callstack = append(m.callstack, "Create")
	model.SetID(m.createMock.ID())

	return m.errMock
}

func (m *MockedTransactionService) ReadAllPurchases(ctx context.Context, id int64) ([]Model, error) {
	m.callstack = append(m.callstack, "ReadAllPurchases")
	return m.readAllPurchasesMock, m.errMock
}

func (m *MockedTransactionService) UpdateAll(ctx context.Context, model []Model) error {
	m.callstack = append(m.callstack, "UpdateAll")
	return m.errMock
}

func (m *MockedTransactionService) Execute(ctx context.Context, op atomic.AtomicOperation[Repository]) error {
	m.callstack = append(m.callstack, "Execute")
	// For test purpose only - In real implementation single and atomic repository should be implemented in distinct structs
	op(ctx, m)

	return m.errMock
}

type MockedAccountService struct {
	readMock *account.Model
	errMock  error
}

func (m *MockedAccountService) Read(ctx context.Context, id int64) (*account.Model, error) {
	return m.readMock, m.errMock
}

func (m *MockedAccountService) Create(ctx context.Context, model *account.Model) error {
	return nil
}

func TestCreateDebtTransaction(t *testing.T) {
	scenarios := []struct {
		given           Input
		expectedAmount  float64
		expectedBalance float64
	}{
		{Input{AccountID: accountID, OperationTypeID: int(Withdrawal), Amount: 24.5}, -24.5, -24.5},
		{Input{AccountID: accountID, OperationTypeID: int(CashPurchase), Amount: 24.5}, -24.5, -24.5},
		{Input{AccountID: accountID, OperationTypeID: int(InstallmentPurchase), Amount: 24.5}, -24.5, -24.5},
	}

	for _, tt := range scenarios {
		tid := int64(100)

		mockTransactRepo := &MockedTransactionService{createMock: &Model{id: tid}}
		mockAccountRepo := &MockedAccountService{readMock: accountDummy}

		uc := NewUsecase(mockTransactRepo, mockAccountRepo)

		ctx := context.Background()
		actual, err := uc.CreateTransaction(ctx, tt.given)

		assert.Nil(t, err)
		assert.Equal(t, tid, actual.id)
		assert.Equal(t, tt.given.AccountID, actual.accountID)
		assert.Equal(t, tt.given.OperationTypeID, int(actual.operationTypeID))
		assert.Equal(t, tt.expectedAmount, actual.amount)
		assert.Equal(t, tt.expectedBalance, actual.balance)
		assert.IsType(t, time.Now(), actual.eventDate)

		spyHasCalled := strings.Join(mockTransactRepo.callstack, "-")
		assert.False(t, strings.Contains(spyHasCalled, "ReadAllPurchases"))
	}
}

func TestCreatePaymentTransaction(t *testing.T) {
	tid := int64(100)
	paymentAmount := float64(200)

	dummyPurchases := []Model{
		*NewTransaction(accountID, int(CashPurchase), 40.0),
		*NewTransaction(accountID, int(InstallmentPurchase), 35.5),
		*NewTransaction(accountID, int(CashPurchase), 15.0),
	}

	mockTransactRepo := &MockedTransactionService{
		createMock:           &Model{id: tid},
		readAllPurchasesMock: dummyPurchases,
	}
	mockAccountRepo := &MockedAccountService{readMock: accountDummy}

	uc := NewUsecase(mockTransactRepo, mockAccountRepo)

	ctx := context.Background()
	input := Input{AccountID: accountID, OperationTypeID: int(Payment), Amount: paymentAmount}

	actual, err := uc.CreateTransaction(ctx, input)

	assert.Nil(t, err)
	assert.Equal(t, tid, actual.id)
	assert.Equal(t, accountID, actual.accountID)
	assert.Equal(t, input.OperationTypeID, int(actual.operationTypeID))
	assert.Equal(t, paymentAmount, actual.amount)
	assert.Equal(t, paymentAmount, actual.balance)
	assert.IsType(t, time.Now(), actual.eventDate)

	spyHasCalled := strings.Join(mockTransactRepo.callstack, "-")
	assert.True(t, strings.Contains(spyHasCalled, "ReadAllPurchases"))

	for _, p := range dummyPurchases {
		assert.Equal(t, 0.0, p.balance)
	}
}

func TestCreatePaymentTransactionWithInsufficientBalanceForAllPurchases(t *testing.T) {
	paymentAmount := float64(50)
	purchases := []Model{
		*NewTransaction(accountID, int(CashPurchase), 40.0),
		*NewTransaction(accountID, int(InstallmentPurchase), 35.5),
		*NewTransaction(accountID, int(CashPurchase), 15.0),
	}

	mockTransactRepo := &MockedTransactionService{
		createMock:           &Model{id: int64(100)},
		readAllPurchasesMock: purchases,
	}

	mockAccountRepo := &MockedAccountService{readMock: accountDummy}

	uc := NewUsecase(mockTransactRepo, mockAccountRepo)

	ctx := context.Background()
	input := Input{AccountID: accountID, OperationTypeID: int(Payment), Amount: paymentAmount}

	_, err := uc.CreateTransaction(ctx, input)
	assert.Nil(t, err)

	spyHasCalled := strings.Join(mockTransactRepo.callstack, "-")
	assert.True(t, strings.Contains(spyHasCalled, "ReadAllPurchases"))

	assert.Equal(t, 0.0, purchases[0].balance)
	assert.Equal(t, -25.5, purchases[1].balance)
	assert.Equal(t, -15.0, purchases[2].balance)
}
