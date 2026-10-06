package ledger

import "errors"

var (
	ErrInvalidID            = errors.New("invalid ledger entry id")
	ErrInvalidWalletID      = errors.New("invalid ledger wallet id")
	ErrInvalidTransactionID = errors.New("invalid ledger transaction id")
	ErrInvalidDirection     = errors.New("invalid ledger direction")
	ErrInvalidAmount        = errors.New("invalid ledger amount")
	ErrInvalidBalance       = errors.New("invalid ledger balance")
	ErrCurrencyMismatch     = errors.New("ledger currency mismatch")
	ErrBalanceMismatch      = errors.New("ledger balance transition mismatch")
	ErrInvalidTimestamp     = errors.New("invalid ledger timestamp")
)
