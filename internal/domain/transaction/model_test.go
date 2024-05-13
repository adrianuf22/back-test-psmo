package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentTransactionAmount(t *testing.T) {
	var dataprovider = []struct {
		given    Model
		expected float64
	}{
		{Model{operationTypeID: Payment, amount: 24.50}, float64(24.50)},  // Optimist
		{Model{operationTypeID: Payment, amount: -24.50}, float64(24.50)}, // Negative value given
		{Model{operationTypeID: Payment, amount: 0}, float64(0)},          // Zeroed
	}

	for _, tt := range dataprovider {
		assert.Equal(t, tt.expected, tt.given.Amount())
	}
}

func TestCashTransactionAmount(t *testing.T) {
	var dataprovider = []struct {
		given    Model
		expected float64
	}{
		{Model{operationTypeID: CashPurchase, amount: -24.50}, float64(-24.50)},       // Optimist
		{Model{operationTypeID: CashPurchase, amount: 24.50}, float64(-24.50)},        // Positive value given
		{Model{operationTypeID: Withdrawal, amount: 24.50}, float64(-24.50)},          // Different type id
		{Model{operationTypeID: InstallmentPurchase, amount: 24.50}, float64(-24.50)}, // Different type id
		{Model{operationTypeID: InstallmentPurchase, amount: 0}, float64(0)},          // Zeroed
	}

	for _, tt := range dataprovider {
		assert.Equal(t, tt.expected, tt.given.Amount())
	}
}
