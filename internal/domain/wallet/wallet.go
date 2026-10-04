package wallet

import (
	"strings"
	"time"

	"github.com/CayoHenri/wagerflow/internal/domain/money"
)

type Wallet struct {
	id        string
	playerID  string
	currency  money.Currency
	balance   money.Money
	version   int64
	createdAt time.Time
	updatedAt time.Time
}

func New(id, playerID string, currency money.Currency, balance money.Money, now time.Time) (*Wallet, error) {
	id = strings.TrimSpace(id)
	playerID = strings.TrimSpace(playerID)

	if id == "" {
		return nil, ErrInvalidID
	}

	if playerID == "" {
		return nil, ErrInvalidPlayerID
	}

	if !currency.IsValid() {
		return nil, money.ErrInvalidCurrency
	}

	if !balance.IsValid() {
		return nil, ErrInvalidBalance
	}

	if balance.Currency() != currency {
		return nil, ErrCurrencyMismatch
	}

	if balance.IsNegative() {
		return nil, ErrInvalidBalance
	}

	return &Wallet{
		id:        id,
		playerID:  playerID,
		currency:  currency,
		balance:   balance,
		version:   1,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Restore(
	id string,
	playerID string,
	currency money.Currency,
	balance money.Money,
	version int64,
	createdAt time.Time,
	updatedAt time.Time,
) (*Wallet, error) {
	id = strings.TrimSpace(id)
	playerID = strings.TrimSpace(playerID)

	if id == "" {
		return nil, ErrInvalidID
	}

	if playerID == "" {
		return nil, ErrInvalidPlayerID
	}

	if !currency.IsValid() {
		return nil, money.ErrInvalidCurrency
	}

	if !balance.IsValid() {
		return nil, ErrInvalidBalance
	}

	if balance.Currency() != currency {
		return nil, ErrCurrencyMismatch
	}

	if balance.IsNegative() {
		return nil, ErrInvalidBalance
	}

	if version < 1 {
		return nil, ErrInvalidVersion
	}

	return &Wallet{
		id:        id,
		playerID:  playerID,
		currency:  currency,
		balance:   balance,
		version:   version,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (w Wallet) ID() string {
	return w.id
}

func (w Wallet) PlayerID() string {
	return w.playerID
}

func (w Wallet) Currency() money.Currency {
	return w.currency
}

func (w Wallet) Balance() money.Money {
	return w.balance
}

func (w Wallet) Version() int64 {
	return w.version
}

func (w Wallet) CreatedAt() time.Time {
	return w.createdAt
}

func (w Wallet) UpdatedAt() time.Time {
	return w.updatedAt
}
