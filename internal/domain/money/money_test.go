package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromMinorUnits(t *testing.T) {
	t.Run("should create money from minor units", func(t *testing.T) {
		m, err := NewFromMinorUnits(1025, BRL)

		require.NoError(t, err)

		assert.Equal(t, int64(1025), m.Amount())
		assert.Equal(t, BRL, m.Currency())
	})

	t.Run("should create zero money", func(t *testing.T) {
		m, err := NewFromMinorUnits(0, BRL)

		require.NoError(t, err)

		assert.Zero(t, m.Amount())
		assert.Equal(t, BRL, m.Currency())
	})

	t.Run("should reject negative minor units", func(t *testing.T) {
		m, err := NewFromMinorUnits(-1, BRL)

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrNegativeAmount)
		assert.Equal(t, Money{}, m)
	})

	t.Run("should reject invalid currency", func(t *testing.T) {
		m, err := NewFromMinorUnits(100, Currency{})

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrInvalidCurrency)
		assert.Equal(t, Money{}, m)
	})
}

func TestParse(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		currency       Currency
		expectedAmount int64
		wantErr        error
	}{
		{
			name:           "should parse integer amount",
			input:          "10",
			currency:       BRL,
			expectedAmount: 1000,
		},
		{
			name:           "should parse amount with one decimal place",
			input:          "10.2",
			currency:       BRL,
			expectedAmount: 1020,
		},
		{
			name:           "should parse amount with two decimal places",
			input:          "10.25",
			currency:       BRL,
			expectedAmount: 1025,
		},
		{
			name:           "should parse zero",
			input:          "0.00",
			currency:       BRL,
			expectedAmount: 0,
		},
		{
			name:           "should parse minimum cent",
			input:          "0.01",
			currency:       BRL,
			expectedAmount: 1,
		},
		{
			name:           "should parse amount with leading zeros",
			input:          "00010.25",
			currency:       BRL,
			expectedAmount: 1025,
		},
		{
			name:           "should trim spaces",
			input:          " 10.25 ",
			currency:       BRL,
			expectedAmount: 1025,
		},
		{
			name:     "should reject empty amount",
			input:    "",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should reject negative amount",
			input:    "-10.00",
			currency: BRL,
			wantErr:  ErrNegativeAmount,
		},
		{
			name:     "should reject more than two decimal places",
			input:    "10.256",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should reject scientific notation",
			input:    "1e3",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should reject NaN",
			input:    "NaN",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should reject Infinity",
			input:    "Infinity",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should reject comma decimal separator",
			input:    "10,25",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should reject trailing decimal separator",
			input:    "10.",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should reject leading decimal separator",
			input:    ".25",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should reject positive sign",
			input:    "+10.00",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should reject invalid text",
			input:    "abc",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
		{
			name:     "should reject multiple decimal separators",
			input:    "10.20.30",
			currency: BRL,
			wantErr:  ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.input, tt.currency)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.expectedAmount, m.Amount())
			assert.Equal(t, tt.currency, m.Currency())
		})
	}
}

func TestParseOverflow(t *testing.T) {
	m, err := Parse("92233720368547758.08", BRL)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrAmountOverflow)
	assert.Equal(t, Money{}, m)
}

func TestParseMaximumAmount(t *testing.T) {
	m, err := Parse("92233720368547758.07", BRL)

	require.NoError(t, err)

	assert.Equal(t, int64(9223372036854775807), m.Amount())
	assert.Equal(t, BRL, m.Currency())
}

func TestZero(t *testing.T) {
	t.Run("should create zero money", func(t *testing.T) {
		m, err := Zero(BRL)

		require.NoError(t, err)

		assert.Zero(t, m.Amount())
		assert.Equal(t, BRL, m.Currency())
		assert.Equal(t, "0.00", m.String())
	})

	t.Run("should reject invalid currency", func(t *testing.T) {
		m, err := Zero(Currency{})

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrInvalidCurrency)
		assert.Equal(t, Money{}, m)
	})
}

func TestMoneyString(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		expected string
	}{
		{
			name:     "should format zero",
			amount:   0,
			expected: "0.00",
		},
		{
			name:     "should format one cent",
			amount:   1,
			expected: "0.01",
		},
		{
			name:     "should format ten cents",
			amount:   10,
			expected: "0.10",
		},
		{
			name:     "should format one real",
			amount:   100,
			expected: "1.00",
		},
		{
			name:     "should format decimal amount",
			amount:   1025,
			expected: "10.25",
		},
		{
			name:     "should format large amount",
			amount:   123456789,
			expected: "1234567.89",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewFromMinorUnits(tt.amount, BRL)

			require.NoError(t, err)

			assert.Equal(t, tt.expected, m.String())
		})
	}
}

func TestMoneyParseAndString(t *testing.T) {
	m, err := Parse("1234567.89", BRL)

	require.NoError(t, err)

	assert.Equal(t, int64(123456789), m.Amount())
	assert.Equal(t, "1234567.89", m.String())
	assert.Equal(t, BRL, m.Currency())
}
