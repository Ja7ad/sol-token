package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToLamport(t *testing.T) {
	tests := []struct {
		input    float64
		expected uint64
	}{
		{1.0, 1_000_000_000},
		{0.5, 500_000_000},
		{0.000000001, 1},
	}

	for _, tt := range tests {
		result := ConvertToLamport(tt.input)
		assert.Equal(t, tt.expected, result, "input: %f", tt.input)
	}
}

func TestConvertToDecimals(t *testing.T) {
	tests := []struct {
		human    float64
		decimals uint8
		expected uint64
	}{
		{1.0, 6, 1000000},
		{0.5, 6, 500000},
		{2.345, 3, 2345},
		{1.23456789, 8, 123456789},
	}

	for _, tt := range tests {
		result := ConvertToDecimals(tt.human, tt.decimals)
		assert.Equal(t, tt.expected, result, "input: %f decimals: %d", tt.human, tt.decimals)
	}
}
