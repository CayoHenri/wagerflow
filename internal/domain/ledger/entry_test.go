package ledger

import (
	"testing"
	"time"

	"github.com/CayoHenri/wagerflow/internal/domain/money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCreditEntry(t *testing.T) {
	now := time.Date(2026, time.October, 6, 12, 0, 0, 0, time.UTC)

	amount, err := money.Parse("25.00", money.BRL)
	require.NoError(t, err)

	before, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	after, err := money.Parse("125.00", money.BRL)
	require.NoError(t, err)

	entry, err := New("ledger-1", "wallet-1", "transaction-1", DirectionCredit, amount, before, after, now)

	require.NoError(t, err)
	require.NotNil(t, entry)

	assert.Equal(t, "ledger-1", entry.ID())
	assert.Equal(t, "wallet-1", entry.WalletID())
	assert.Equal(t, "transaction-1", entry.TransactionID())
	assert.Equal(t, DirectionCredit, entry.Direction())
	assert.True(t, amount.Equal(entry.Amount()))
	assert.True(t, before.Equal(entry.BalanceBefore()))
	assert.True(t, after.Equal(entry.BalanceAfter()))
	assert.Equal(t, now, entry.CreatedAt())
}

func TestNewDebitEntry(t *testing.T) {
	now := time.Now().UTC()

	amount, err := money.Parse("25.00", money.BRL)
	require.NoError(t, err)

	before, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	after, err := money.Parse("75.00", money.BRL)
	require.NoError(t, err)

	entry, err := New("ledger-1", "wallet-1", "transaction-1", DirectionDebit, amount, before, after, now)

	require.NoError(t, err)
	require.NotNil(t, entry)

	assert.Equal(t, DirectionDebit, entry.Direction())
	assert.Equal(t, "25.00", entry.Amount().String())
	assert.Equal(t, "100.00", entry.BalanceBefore().String())
	assert.Equal(t, "75.00", entry.BalanceAfter().String())
}

func TestNewRejectsInvalidBalanceTransition(t *testing.T) {
	amount, err := money.Parse("20.00", money.BRL)
	require.NoError(t, err)

	before, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	after, err := money.Parse("90.00", money.BRL)
	require.NoError(t, err)

	entry, err := New("ledger-1", "wallet-1", "transaction-1", DirectionDebit, amount, before, after, time.Now().UTC())

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrBalanceMismatch)
	assert.Nil(t, entry)
}

func TestNewRejectsCurrencyMismatch(t *testing.T) {
	amount, err := money.Parse("20.00", money.BRL)
	require.NoError(t, err)

	before, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	after, err := money.Parse("120.00", money.USD)
	require.NoError(t, err)

	entry, err := New("ledger-1", "wallet-1", "transaction-1", DirectionCredit, amount, before, after, time.Now().UTC())

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrCurrencyMismatch)
	assert.Nil(t, entry)
}

func TestNewRejectsNonPositiveAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
	}{
		{
			name:   "zero amount",
			amount: 0,
		},
		{
			name:   "negative amount",
			amount: -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount, err := money.NewFromMinorUnits(tt.amount, money.BRL)
			require.NoError(t, err)

			before, err := money.Parse("100.00", money.BRL)
			require.NoError(t, err)

			after, err := money.Parse("100.00", money.BRL)
			require.NoError(t, err)

			entry, err := New("ledger-1", "wallet-1", "transaction-1", DirectionCredit, amount, before, after, time.Now().UTC())

			require.Error(t, err)

			assert.ErrorIs(t, err, ErrInvalidAmount)
			assert.Nil(t, entry)
		})
	}
}

func TestNewRejectsInvalidIdentifiers(t *testing.T) {
	amount, err := money.Parse("10.00", money.BRL)
	require.NoError(t, err)

	before, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	after, err := money.Parse("110.00", money.BRL)
	require.NoError(t, err)

	tests := []struct {
		name          string
		id            string
		walletID      string
		transactionID string
		expectedErr   error
	}{
		{
			name:          "empty id",
			id:            "",
			walletID:      "wallet-1",
			transactionID: "transaction-1",
			expectedErr:   ErrInvalidID,
		},
		{
			name:          "empty wallet id",
			id:            "ledger-1",
			walletID:      "",
			transactionID: "transaction-1",
			expectedErr:   ErrInvalidWalletID,
		},
		{
			name:          "empty transaction id",
			id:            "ledger-1",
			walletID:      "wallet-1",
			transactionID: "",
			expectedErr:   ErrInvalidTransactionID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := New(tt.id, tt.walletID, tt.transactionID, DirectionCredit, amount, before, after, time.Now().UTC())

			require.Error(t, err)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Nil(t, entry)
		})
	}
}

func TestNewRejectsInvalidDirection(t *testing.T) {
	amount, err := money.Parse("10.00", money.BRL)
	require.NoError(t, err)

	before, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	after, err := money.Parse("110.00", money.BRL)
	require.NoError(t, err)

	entry, err := New("ledger-1", "wallet-1", "transaction-1", Direction("INVALID"), amount, before, after, time.Now().UTC())

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInvalidDirection)
	assert.Nil(t, entry)
}

func TestNewRejectsZeroTimestamp(t *testing.T) {
	amount, err := money.Parse("10.00", money.BRL)
	require.NoError(t, err)

	before, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	after, err := money.Parse("110.00", money.BRL)
	require.NoError(t, err)

	entry, err := New("ledger-1", "wallet-1", "transaction-1", DirectionCredit, amount, before, after, time.Time{})

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInvalidTimestamp)
	assert.Nil(t, entry)
}
