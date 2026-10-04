package money

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustMoneyFromMinorUnits(t *testing.T, amount int64, currency Currency) Money {
	t.Helper()

	m, err := NewFromMinorUnits(amount, currency)
	require.NoError(t, err)

	return m
}

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

	t.Run("should create negative money", func(t *testing.T) {
		m, err := NewFromMinorUnits(-10, BRL)
		require.NoError(t, err)

		assert.Equal(t, int64(-10), m.Amount())
		assert.Equal(t, BRL, m.Currency())
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

func TestMoneyIsZero(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		expected bool
	}{
		{
			name:     "should return true when amount is zero",
			amount:   0,
			expected: true,
		},
		{
			name:     "should return false when amount is greater than zero",
			amount:   1,
			expected: false,
		},
		{
			name:     "should return false for positive amount",
			amount:   1025,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewFromMinorUnits(tt.amount, BRL)

			require.NoError(t, err)

			assert.Equal(t, tt.expected, m.IsZero())
		})
	}
}

func TestMoneyEqual(t *testing.T) {
	tests := []struct {
		name     string
		first    Money
		second   Money
		expected bool
	}{
		{
			name:     "should return true for equal money",
			first:    mustMoneyFromMinorUnits(t, 1000, BRL),
			second:   mustMoneyFromMinorUnits(t, 1000, BRL),
			expected: true,
		},
		{
			name:     "should return false for different amounts",
			first:    mustMoneyFromMinorUnits(t, 1000, BRL),
			second:   mustMoneyFromMinorUnits(t, 2000, BRL),
			expected: false,
		},
		{
			name:     "should return false for different currencies",
			first:    mustMoneyFromMinorUnits(t, 1000, BRL),
			second:   mustMoneyFromMinorUnits(t, 1000, USD),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.first.Equal(tt.second))
		})
	}
}

func TestMoneyCompare(t *testing.T) {
	tests := []struct {
		name     string
		first    Money
		second   Money
		expected int
		wantErr  error
	}{
		{
			name:     "should return minus one when first amount is smaller",
			first:    mustMoneyFromMinorUnits(t, 1000, BRL),
			second:   mustMoneyFromMinorUnits(t, 2000, BRL),
			expected: -1,
		},
		{
			name:     "should return zero when amounts are equal",
			first:    mustMoneyFromMinorUnits(t, 1000, BRL),
			second:   mustMoneyFromMinorUnits(t, 1000, BRL),
			expected: 0,
		},
		{
			name:     "should return one when first amount is greater",
			first:    mustMoneyFromMinorUnits(t, 2000, BRL),
			second:   mustMoneyFromMinorUnits(t, 1000, BRL),
			expected: 1,
		},
		{
			name:    "should reject different currencies",
			first:   mustMoneyFromMinorUnits(t, 1000, BRL),
			second:  mustMoneyFromMinorUnits(t, 1000, USD),
			wantErr: ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.first.Compare(tt.second)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMoneyIsValid(t *testing.T) {
	t.Run("should return true for valid money", func(t *testing.T) {
		m, err := Parse("10.25", BRL)

		require.NoError(t, err)

		assert.True(t, m.IsValid())
	})

	t.Run("should return true for valid zero money", func(t *testing.T) {
		m, err := Zero(BRL)

		require.NoError(t, err)

		assert.True(t, m.IsValid())
	})

	t.Run("should return false for zero value money", func(t *testing.T) {
		var m Money

		assert.False(t, m.IsValid())
	})
}

func TestMoneyAdd(t *testing.T) {
	t.Run("should add money with same currency", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		second, err := Parse("25.50", BRL)
		require.NoError(t, err)

		result, err := first.Add(second)

		require.NoError(t, err)

		assert.Equal(t, "125.50", result.String())
		assert.Equal(t, BRL, result.Currency())
	})

	t.Run("should add zero", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		zero, err := Zero(BRL)
		require.NoError(t, err)

		result, err := first.Add(zero)

		require.NoError(t, err)

		assert.True(t, first.Equal(result))
	})

	t.Run("should reject different currencies", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		second, err := Parse("25.00", USD)
		require.NoError(t, err)

		result, err := first.Add(second)

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrCurrencyMismatch)
		assert.Equal(t, Money{}, result)
	})

	t.Run("should reject overflow", func(t *testing.T) {
		first, err := NewFromMinorUnits(math.MaxInt64, BRL)
		require.NoError(t, err)

		second, err := NewFromMinorUnits(1, BRL)
		require.NoError(t, err)

		result, err := first.Add(second)

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrAmountOverflow)
		assert.Equal(t, Money{}, result)
	})

	t.Run("should reject invalid money", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		var invalid Money

		result, err := first.Add(invalid)

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrInvalidMoney)
		assert.Equal(t, Money{}, result)
	})
}

func TestMoneySubtract(t *testing.T) {
	t.Run("should subtract money with same currency", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		second, err := Parse("25.50", BRL)
		require.NoError(t, err)

		result, err := first.Subtract(second)

		require.NoError(t, err)

		assert.Equal(t, "74.50", result.String())
		assert.Equal(t, BRL, result.Currency())
	})

	t.Run("should result in zero when amounts are equal", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		second, err := Parse("100.00", BRL)
		require.NoError(t, err)

		result, err := first.Subtract(second)

		require.NoError(t, err)

		assert.True(t, result.IsZero())
		assert.Equal(t, "0.00", result.String())
	})

	t.Run("should subtract zero", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		zero, err := Zero(BRL)
		require.NoError(t, err)

		result, err := first.Subtract(zero)

		require.NoError(t, err)

		assert.True(t, first.Equal(result))
	})

	t.Run("should allow negative subtraction result", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		second, err := Parse("150.00", BRL)
		require.NoError(t, err)

		result, err := first.Subtract(second)

		require.NoError(t, err)

		assert.Equal(t, int64(-5000), result.Amount())
		assert.Equal(t, "-50.00", result.String())
		assert.True(t, result.IsNegative())
	})

	t.Run("should reject different currencies", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		second, err := Parse("25.00", USD)
		require.NoError(t, err)

		result, err := first.Subtract(second)

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrCurrencyMismatch)
		assert.Equal(t, Money{}, result)
	})

	t.Run("should reject invalid money", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		var invalid Money

		result, err := first.Subtract(invalid)

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrInvalidMoney)
		assert.Equal(t, Money{}, result)
	})

	t.Run("should add negative money", func(t *testing.T) {
		first, err := Parse("100.00", BRL)
		require.NoError(t, err)

		second, err := NewFromMinorUnits(-2500, BRL)
		require.NoError(t, err)

		result, err := first.Add(second)

		require.NoError(t, err)

		assert.Equal(t, "75.00", result.String())
	})

	t.Run("should reject underflow", func(t *testing.T) {
		first, err := NewFromMinorUnits(math.MinInt64, BRL)
		require.NoError(t, err)

		second, err := NewFromMinorUnits(-1, BRL)
		require.NoError(t, err)

		result, err := first.Add(second)

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrAmountOverflow)
		assert.Equal(t, Money{}, result)
	})
}

func TestMoneySign(t *testing.T) {
	tests := []struct {
		name       string
		amount     int64
		isNegative bool
		isZero     bool
		isPositive bool
	}{
		{
			name:       "should identify negative money",
			amount:     -100,
			isNegative: true,
		},
		{
			name:   "should identify zero money",
			amount: 0,
			isZero: true,
		},
		{
			name:       "should identify positive money",
			amount:     100,
			isPositive: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewFromMinorUnits(tt.amount, BRL)

			require.NoError(t, err)

			assert.Equal(t, tt.isNegative, m.IsNegative())
			assert.Equal(t, tt.isZero, m.IsZero())
			assert.Equal(t, tt.isPositive, m.IsPositive())
		})
	}
}

func TestMoneyNegate(t *testing.T) {
	t.Run("should negate positive money", func(t *testing.T) {
		m, err := Parse("10.25", BRL)
		require.NoError(t, err)

		result, err := m.Negate()

		require.NoError(t, err)

		assert.Equal(t, int64(-1025), result.Amount())
		assert.Equal(t, "-10.25", result.String())
		assert.Equal(t, BRL, result.Currency())
		assert.True(t, result.IsNegative())
	})

	t.Run("should negate negative money", func(t *testing.T) {
		m, err := NewFromMinorUnits(-1025, BRL)
		require.NoError(t, err)

		result, err := m.Negate()

		require.NoError(t, err)

		assert.Equal(t, int64(1025), result.Amount())
		assert.Equal(t, "10.25", result.String())
		assert.True(t, result.IsPositive())
	})

	t.Run("should keep zero as zero", func(t *testing.T) {
		m, err := Zero(BRL)
		require.NoError(t, err)

		result, err := m.Negate()

		require.NoError(t, err)

		assert.True(t, result.IsZero())
		assert.Equal(t, "0.00", result.String())
	})

	t.Run("should reject minimum int64 overflow", func(t *testing.T) {
		m, err := NewFromMinorUnits(math.MinInt64, BRL)
		require.NoError(t, err)

		result, err := m.Negate()

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrAmountOverflow)
		assert.Equal(t, Money{}, result)
	})

	t.Run("should reject invalid money", func(t *testing.T) {
		var m Money

		result, err := m.Negate()

		require.Error(t, err)

		assert.ErrorIs(t, err, ErrInvalidMoney)
		assert.Equal(t, Money{}, result)
	})
}
