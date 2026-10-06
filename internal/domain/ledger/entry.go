package ledger

import (
	"strings"
	"time"

	"github.com/CayoHenri/wagerflow/internal/domain/money"
)

type WalletLedgerEntry struct {
	id            string
	walletID      string
	transactionID string
	direction     Direction
	amount        money.Money
	balanceBefore money.Money
	balanceAfter  money.Money
	createdAt     time.Time
}

func New(
	id string,
	walletID string,
	transactionID string,
	direction Direction,
	amount money.Money,
	balanceBefore money.Money,
	balanceAfter money.Money,
	createdAt time.Time,
) (*WalletLedgerEntry, error) {
	id = strings.TrimSpace(id)
	walletID = strings.TrimSpace(walletID)
	transactionID = strings.TrimSpace(transactionID)

	if id == "" {
		return nil, ErrInvalidID
	}

	if walletID == "" {
		return nil, ErrInvalidWalletID
	}

	if transactionID == "" {
		return nil, ErrInvalidTransactionID
	}

	if !direction.IsValid() {
		return nil, ErrInvalidDirection
	}

	if !amount.IsValid() || !amount.IsPositive() {
		return nil, ErrInvalidAmount
	}

	if !balanceBefore.IsValid() || !balanceAfter.IsValid() || balanceBefore.IsNegative() || balanceAfter.IsNegative() {
		return nil, ErrInvalidBalance
	}

	if createdAt.IsZero() {
		return nil, ErrInvalidTimestamp
	}

	if amount.Currency() != balanceBefore.Currency() || amount.Currency() != balanceAfter.Currency() {
		return nil, ErrCurrencyMismatch
	}

	if err := validateTransition(
		direction,
		amount,
		balanceBefore,
		balanceAfter,
	); err != nil {
		return nil, err
	}

	return &WalletLedgerEntry{
		id:            id,
		walletID:      walletID,
		transactionID: transactionID,
		direction:     direction,
		amount:        amount,
		balanceBefore: balanceBefore,
		balanceAfter:  balanceAfter,
		createdAt:     createdAt,
	}, nil
}

func (e WalletLedgerEntry) ID() string {
	return e.id
}

func (e WalletLedgerEntry) WalletID() string {
	return e.walletID
}

func (e WalletLedgerEntry) TransactionID() string {
	return e.transactionID
}

func (e WalletLedgerEntry) Direction() Direction {
	return e.direction
}

func (e WalletLedgerEntry) Amount() money.Money {
	return e.amount
}

func (e WalletLedgerEntry) BalanceBefore() money.Money {
	return e.balanceBefore
}

func (e WalletLedgerEntry) BalanceAfter() money.Money {
	return e.balanceAfter
}

func (e WalletLedgerEntry) CreatedAt() time.Time {
	return e.createdAt
}

func validateTransition(direction Direction, amount money.Money, balanceBefore money.Money, balanceAfter money.Money) error {
	var expected money.Money
	var err error

	switch direction {
	case DirectionCredit:
		expected, err = balanceBefore.Add(amount)

	case DirectionDebit:
		expected, err = balanceBefore.Subtract(amount)

	default:
		return ErrInvalidDirection
	}

	if err != nil {
		return err
	}

	if expected.IsNegative() {
		return ErrInvalidBalance
	}

	if !expected.Equal(balanceAfter) {
		return ErrBalanceMismatch
	}

	return nil
}
