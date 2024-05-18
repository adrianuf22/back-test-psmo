package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentTransactionAmount(t *testing.T) {
	dataprovider := []struct {
		given    Model
		expected float64
	}{
		{*NewTransaction(1, int(Payment), 24.50), float64(24.50)},  // Optimist
		{*NewTransaction(1, int(Payment), -24.50), float64(24.50)}, // Negative value given
		{*NewTransaction(1, int(Payment), 0), float64(0)},          // Zeroed
	}

	for _, tt := range dataprovider {
		assert.Equal(t, tt.expected, tt.given.Amount())
	}
}

func TestCashTransactionAmount(t *testing.T) {
	dataprovider := []struct {
		given    Model
		expected float64
	}{
		{*NewTransaction(1, int(CashPurchase), -24.50), float64(-24.50)},       // Optimist
		{*NewTransaction(1, int(CashPurchase), 24.50), float64(-24.50)},        // Positive value given
		{*NewTransaction(1, int(Withdrawal), 24.50), float64(-24.50)},          // Different type id
		{*NewTransaction(1, int(InstallmentPurchase), 24.50), float64(-24.50)}, // Different type id
		{*NewTransaction(1, int(InstallmentPurchase), 0), float64(0)},          // Zeroed
	}

	for _, tt := range dataprovider {
		assert.Equal(t, tt.expected, tt.given.Amount())
	}
}

func TestPaymentDischarge(t *testing.T) {
	dataprovider := []struct {
		given           Model
		input           float64
		expectedBalance float64
		expectedOutput  float64
	}{
		{Model{amount: -60.0, balance: -60.0}, 100.0, 0, 40.0}, // Input is enough - Full payment
		{Model{amount: -60.0, balance: -60.0}, 50.0, -10.0, 0}, // Input is lower than amount - Partial payment
		{Model{amount: -40.0, balance: -13.0}, 15.0, 0, 2.0},   // Input cover the rest before payed balance - Partial payment
	}

	for _, tt := range dataprovider {
		output, err := tt.given.Discharge(tt.input)

		assert.Nil(t, err)
		assert.Equal(t, tt.expectedBalance, tt.given.Balance())
		assert.Equal(t, tt.expectedOutput, output)
	}
}

func TestPaymentDischargeFailsWhenTryToDischargeAPaymentTransaction(t *testing.T) {
	m := Model{amount: 60.0, balance: 60.0, operationTypeID: Payment}

	a, err := m.Discharge(100)

	assert.EqualError(t, err, ErrOperationTypeNotAllowed.Error())
	assert.Equal(t, 100.0, a)
}

func TestPaymentDischargeFailsWhenTryToDischargeWithZeroAsValue(t *testing.T) {
	m := Model{amount: 60.0, balance: 60.0, operationTypeID: CashPurchase}

	a, err := m.Discharge(0)

	assert.EqualError(t, err, ErrInsufficientAmount.Error())
	assert.Equal(t, 0.0, a)
}
