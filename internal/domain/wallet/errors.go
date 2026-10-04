package wallet

import "errors"

var (
	ErrInvalidID        = errors.New("invalid wallet id")
	ErrInvalidPlayerID  = errors.New("invalid player id")
	ErrInvalidBalance   = errors.New("invalid wallet balance")
	ErrCurrencyMismatch = errors.New("wallet currency mismatch")
	ErrInvalidVersion   = errors.New("invalid wallet version")
)
