package money

import (
	"strings"
)

type Currency struct {
	code string
}

var (
	BRL = Currency{code: "BRL"}
	USD = Currency{code: "USD"}
	EUR = Currency{code: "EUR"}
)

func NewCurrency(value string) (Currency, error) {
	value = strings.TrimSpace(value)

	switch value {
	case BRL.code:
		return BRL, nil

	case USD.code:
		return USD, nil

	case EUR.code:
		return EUR, nil

	default:
		return Currency{}, ErrInvalidCurrency
	}
}

func (c Currency) String() string {
	return c.code
}

func (c Currency) IsValid() bool {
	switch c {
	case BRL, USD, EUR:
		return true
	default:
		return false
	}
}
