package money

import "errors"

var (
	ErrInvalidAmount    = errors.New("invalid money amount")
	ErrInvalidCurrency  = errors.New("invalid currency")
	ErrNegativeAmount   = errors.New("money amount cannot be negative")
	ErrAmountOverflow   = errors.New("money amount overflow")
	ErrCurrencyMismatch = errors.New("currency mismatch")
	ErrInvalidMoney     = errors.New("invalid money")
)
