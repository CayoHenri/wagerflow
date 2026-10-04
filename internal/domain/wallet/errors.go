package wallet

import "errors"

var (
	ErrInvalidID         = errors.New("invalid wallet id")
	ErrInvalidPlayerID   = errors.New("invalid player id")
	ErrInvalidBalance    = errors.New("invalid wallet balance")
	ErrInvalidVersion    = errors.New("invalid wallet version")
	ErrInvalidAmount     = errors.New("wallet amount must be greater than zero")
	ErrCurrencyMismatch  = errors.New("wallet currency mismatch")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidTimestamp  = errors.New("invalid wallet timestamp")
)
