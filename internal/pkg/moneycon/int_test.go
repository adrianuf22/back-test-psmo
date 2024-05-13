package moneycon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatToInt(t *testing.T) {
	dataprovider := []struct {
		given    int64
		expected int64
	}{
		{FloatToInt(0.05), int64(5)},
		{FloatToInt(3.75), int64(375)},
		{FloatToInt(24.5), int64(2450)},
		{FloatToInt(10350.98), int64(1035098)},
	}

	for _, tt := range dataprovider {
		assert.Equal(t, tt.expected, tt.given)
	}
}

func TestIntToFloat(t *testing.T) {
	dataprovider := []struct {
		given    float64
		expected float64
	}{
		{IntToFloat(5), 0.05},
		{IntToFloat(375), 3.75},
		{IntToFloat(2450), 24.5},
		{IntToFloat(1035098), 10350.98},
	}

	for _, tt := range dataprovider {
		assert.Equal(t, tt.expected, tt.given)
	}
}
