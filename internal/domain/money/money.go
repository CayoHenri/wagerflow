package money

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Money struct {
	amount   int64
	currency Currency
}

// Criação
func NewFromMinorUnits(amount int64, currency Currency) (Money, error) {
	if !currency.IsValid() {
		return Money{}, ErrInvalidCurrency
	}

	if amount < 0 {
		return Money{}, ErrNegativeAmount
	}

	return Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func Parse(value string, currency Currency) (Money, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Money{}, ErrInvalidAmount
	}

	if strings.HasPrefix(value, "-") {
		return Money{}, ErrNegativeAmount
	}

	parts := strings.Split(value, ".")

	if len(parts) > 2 {
		return Money{}, ErrInvalidAmount
	}

	wholePart := parts[0]

	if !isDigits(wholePart) {
		return Money{}, ErrInvalidAmount
	}

	fractionPart := ""

	if len(parts) == 2 {
		fractionPart = parts[1]

		if fractionPart == "" {
			return Money{}, ErrInvalidAmount
		}

		if !isDigits(fractionPart) {
			return Money{}, ErrInvalidAmount
		}
	}

	if len(fractionPart) > 2 {
		return Money{}, ErrInvalidAmount
	}

	whole, err := strconv.ParseInt(wholePart, 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("%w: %s", ErrInvalidAmount, value)
	}

	fraction := int64(0)

	if fractionPart != "" {
		if len(fractionPart) == 1 {
			fractionPart += "0"
		}

		fraction, err = strconv.ParseInt(fractionPart, 10, 64)
		if err != nil {
			return Money{}, fmt.Errorf("%w: %s", ErrInvalidAmount, value)
		}
	}

	if whole > (math.MaxInt64-fraction)/100 {
		return Money{}, ErrAmountOverflow
	}

	amount := whole*100 + fraction

	return NewFromMinorUnits(amount, currency)
}

func Zero(currency Currency) (Money, error) {
	return NewFromMinorUnits(0, currency)
}

// ações
func (m Money) Add(other Money) (Money, error) {
	if !m.IsValid() || !other.IsValid() {
		return Money{}, ErrInvalidMoney
	}

	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	if other.amount > math.MaxInt64-m.amount {
		return Money{}, ErrAmountOverflow
	}

	return NewFromMinorUnits(m.amount+other.amount, m.currency)
}

func (m Money) Subtract(other Money) (Money, error) {
	if !m.IsValid() || !other.IsValid() {
		return Money{}, ErrInvalidMoney
	}

	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	if other.amount > m.amount {
		return Money{}, ErrNegativeAmount
	}

	return NewFromMinorUnits(m.amount-other.amount, m.currency)
}

// representação
func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() Currency {
	return m.currency
}

func (m Money) String() string {
	whole := m.amount / 100
	fraction := m.amount % 100

	return fmt.Sprintf("%d.%02d", whole, fraction)
}

// estado
func (m Money) IsZero() bool {
	return m.IsValid() && m.amount == 0
}

func (m Money) IsValid() bool {
	return m.currency.IsValid() && m.amount >= 0
}

// comparação
func (m Money) Equal(other Money) bool {
	if !m.IsValid() || !other.IsValid() {
		return false
	}

	return m.amount == other.amount && m.currency == other.currency
}

func (m Money) Compare(other Money) (int, error) {
	if !m.IsValid() || !other.IsValid() {
		return 0, ErrInvalidMoney
	}

	if m.currency != other.currency {
		return 0, ErrCurrencyMismatch
	}

	switch {
	case m.amount < other.amount:
		return -1, nil

	case m.amount > other.amount:
		return 1, nil

	default:
		return 0, nil
	}
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}

	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}
