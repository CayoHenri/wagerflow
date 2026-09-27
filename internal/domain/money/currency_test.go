package money

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCurrency(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Currency
		wantErr  error
	}{
		{
			name:     "should create BRL currency",
			input:    "BRL",
			expected: BRL,
		},
		{
			name:     "should create USD currency",
			input:    "USD",
			expected: USD,
		},
		{
			name:     "should create EUR currency",
			input:    "EUR",
			expected: EUR,
		},
		{
			name:     "should trim spaces",
			input:    " BRL ",
			expected: BRL,
		},
		{
			name:    "should reject empty currency",
			input:   "",
			wantErr: ErrInvalidCurrency,
		},
		{
			name:    "should reject unsupported currency",
			input:   "ABC",
			wantErr: ErrInvalidCurrency,
		},
		{
			name:    "should reject lowercase currency",
			input:   "brl",
			wantErr: ErrInvalidCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currency, err := NewCurrency(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, currency)
		})
	}
}

func TestCurrencyString(t *testing.T) {
	currency := BRL

	assert.Equal(t, "BRL", currency.String())
}
