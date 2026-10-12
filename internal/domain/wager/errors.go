package wager

import "errors"

var (
	ErrInvalidID                    = errors.New("invalid wager transaction id")
	ErrInvalidExternalTransactionID = errors.New("invalid external transaction id")
	ErrInvalidProviderID            = errors.New("invalid provider id")
	ErrInvalidIdempotencyKey        = errors.New("invalid idempotency key")
	ErrInvalidPayloadHash           = errors.New("invalid payload hash")
	ErrInvalidWalletID              = errors.New("invalid wallet id")
	ErrInvalidPlayerID              = errors.New("invalid player id")
	ErrInvalidRoundID               = errors.New("invalid round id")
	ErrInvalidGameID                = errors.New("invalid game id")
	ErrInvalidKind                  = errors.New("invalid wager transaction kind")
	ErrInvalidState                 = errors.New("invalid wager transaction state")
	ErrInvalidAmount                = errors.New("invalid wager transaction amount")
	ErrInvalidTimestamp             = errors.New("invalid wager transaction timestamp")
)
