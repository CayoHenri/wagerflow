package wallet

import (
	"testing"
	"time"

	"github.com/CayoHenri/wagerflow/internal/domain/money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("should create wallet", func(t *testing.T) {
		now := time.Date(2026, time.October, 4, 10, 0, 0, 0, time.UTC)

		balance, err := money.Parse("100.00", money.BRL)
		require.NoError(t, err)

		w, err := New("wallet-1", "player-1", money.BRL, balance, now)

		require.NoError(t, err)
		require.NotNil(t, w)

		assert.Equal(t, "wallet-1", w.ID())
		assert.Equal(t, "player-1", w.PlayerID())
		assert.Equal(t, money.BRL, w.Currency())
		assert.True(t, balance.Equal(w.Balance()))
		assert.Equal(t, int64(1), w.Version())
		assert.Equal(t, now, w.CreatedAt())
		assert.Equal(t, now, w.UpdatedAt())
	})
}

func TestNewInvalidTimestamp(t *testing.T) {
	balance, err := money.Zero(money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "player-1", money.BRL, balance, time.Time{})

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidTimestamp)
	assert.Nil(t, w)
}

func TestNewBlankPlayerID(t *testing.T) {
	balance, err := money.Zero(money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "  ", money.BRL, balance, time.Now())

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidPlayerID)
	assert.Nil(t, w)
}

func TestNewInvalidCurrency(t *testing.T) {
	balance, err := money.Zero(money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "player-1", money.Currency{}, balance, time.Now())

	require.Error(t, err)
	assert.ErrorIs(t, err, money.ErrInvalidCurrency)
	assert.Nil(t, w)
}

func TestNewInvalidBalance(t *testing.T) {
	w, err := New("wallet-1", "player-1", money.BRL, money.Money{}, time.Now())

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidBalance)
	assert.Nil(t, w)
}

func TestNewInvalidID(t *testing.T) {
	balance, err := money.Zero(money.BRL)
	require.NoError(t, err)

	w, err := New("", "player-1", money.BRL, balance, time.Now())

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInvalidID)
	assert.Nil(t, w)
}

func TestNewBlankID(t *testing.T) {
	balance, err := money.Zero(money.BRL)
	require.NoError(t, err)

	w, err := New("   ", "player-1", money.BRL, balance, time.Now())

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInvalidID)
	assert.Nil(t, w)
}

func TestNewInvalidPlayerID(t *testing.T) {
	balance, err := money.Zero(money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "", money.BRL, balance, time.Now())

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInvalidPlayerID)
	assert.Nil(t, w)
}

func TestNewCurrencyMismatch(t *testing.T) {
	balance, err := money.Parse("100.00", money.USD)
	require.NoError(t, err)

	w, err := New("wallet-1", "player-1", money.BRL, balance, time.Now())

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrCurrencyMismatch)
	assert.Nil(t, w)
}

func TestNewNegativeBalance(t *testing.T) {
	balance, err := money.NewFromMinorUnits(-1000, money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "player-1", money.BRL, balance, time.Now())

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInvalidBalance)
	assert.Nil(t, w)
}

func TestNewZeroBalance(t *testing.T) {
	balance, err := money.Zero(money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "player-1", money.BRL, balance, time.Now())

	require.NoError(t, err)
	require.NotNil(t, w)

	assert.True(t, w.Balance().IsZero())
	assert.Equal(t, int64(1), w.Version())
}

func TestRestore(t *testing.T) {
	createdAt := time.Date(2026, time.October, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, time.October, 4, 10, 0, 0, 0, time.UTC)

	balance, err := money.Parse("153.42", money.BRL)
	require.NoError(t, err)

	w, err := Restore("wallet-123", "player-456", money.BRL, balance, 17, createdAt, updatedAt)

	require.NoError(t, err)
	require.NotNil(t, w)

	assert.Equal(t, "wallet-123", w.ID())
	assert.Equal(t, "player-456", w.PlayerID())
	assert.Equal(t, money.BRL, w.Currency())
	assert.Equal(t, "153.42", w.Balance().String())
	assert.Equal(t, int64(17), w.Version())
	assert.Equal(t, createdAt, w.CreatedAt())
	assert.Equal(t, updatedAt, w.UpdatedAt())
}

func TestRestoreInvalidState(t *testing.T) {
	createdAt := time.Date(2026, time.October, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	balance, err := money.Zero(money.BRL)
	require.NoError(t, err)
	invalidBalance := money.Money{}
	negativeBalance, err := money.NewFromMinorUnits(-1, money.BRL)
	require.NoError(t, err)
	usdBalance, err := money.Zero(money.USD)
	require.NoError(t, err)

	tests := []struct {
		name      string
		id        string
		playerID  string
		currency  money.Currency
		balance   money.Money
		version   int64
		createdAt time.Time
		updatedAt time.Time
		wantErr   error
	}{
		{name: "empty created timestamp", id: "wallet-1", playerID: "player-1", currency: money.BRL, balance: balance, version: 1, createdAt: time.Time{}, updatedAt: updatedAt, wantErr: ErrInvalidTimestamp},
		{name: "empty updated timestamp", id: "wallet-1", playerID: "player-1", currency: money.BRL, balance: balance, version: 1, createdAt: createdAt, updatedAt: time.Time{}, wantErr: ErrInvalidTimestamp},
		{name: "updated timestamp before creation", id: "wallet-1", playerID: "player-1", currency: money.BRL, balance: balance, version: 1, createdAt: updatedAt, updatedAt: createdAt, wantErr: ErrInvalidTimestamp},
		{name: "empty ID", id: " ", playerID: "player-1", currency: money.BRL, balance: balance, version: 1, createdAt: createdAt, updatedAt: updatedAt, wantErr: ErrInvalidID},
		{name: "empty player ID", id: "wallet-1", playerID: " ", currency: money.BRL, balance: balance, version: 1, createdAt: createdAt, updatedAt: updatedAt, wantErr: ErrInvalidPlayerID},
		{name: "invalid currency", id: "wallet-1", playerID: "player-1", currency: money.Currency{}, balance: balance, version: 1, createdAt: createdAt, updatedAt: updatedAt, wantErr: money.ErrInvalidCurrency},
		{name: "invalid balance", id: "wallet-1", playerID: "player-1", currency: money.BRL, balance: invalidBalance, version: 1, createdAt: createdAt, updatedAt: updatedAt, wantErr: ErrInvalidBalance},
		{name: "currency mismatch", id: "wallet-1", playerID: "player-1", currency: money.BRL, balance: usdBalance, version: 1, createdAt: createdAt, updatedAt: updatedAt, wantErr: ErrCurrencyMismatch},
		{name: "negative balance", id: "wallet-1", playerID: "player-1", currency: money.BRL, balance: negativeBalance, version: 1, createdAt: createdAt, updatedAt: updatedAt, wantErr: ErrInvalidBalance},
		{name: "zero version", id: "wallet-1", playerID: "player-1", currency: money.BRL, balance: balance, version: 0, createdAt: createdAt, updatedAt: updatedAt, wantErr: ErrInvalidVersion},
		{name: "negative version", id: "wallet-1", playerID: "player-1", currency: money.BRL, balance: balance, version: -1, createdAt: createdAt, updatedAt: updatedAt, wantErr: ErrInvalidVersion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := Restore(tt.id, tt.playerID, tt.currency, tt.balance, tt.version, tt.createdAt, tt.updatedAt)

			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, w)
		})
	}
}

func TestWalletCredit(t *testing.T) {
	t.Run("should credit wallet", func(t *testing.T) {
		createdAt := time.Date(2026, time.October, 4, 10, 0, 0, 0, time.UTC)
		creditedAt := createdAt.Add(time.Hour)

		balance, err := money.Parse("100.00", money.BRL)
		require.NoError(t, err)

		w, err := New("wallet-1", "player-1", money.BRL, balance, createdAt)
		require.NoError(t, err)

		amount, err := money.Parse("25.50", money.BRL)
		require.NoError(t, err)

		err = w.Credit(amount, creditedAt)

		require.NoError(t, err)

		assert.Equal(t, "125.50", w.Balance().String())
		assert.Equal(t, int64(2), w.Version())
		assert.Equal(t, creditedAt, w.UpdatedAt())

		// A criação da Wallet não mudou.
		assert.Equal(t, createdAt, w.CreatedAt())
	})
}

func TestWalletCreditZeroAmount(t *testing.T) {
	now := time.Now()

	balance, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "player-1", money.BRL, balance, now)
	require.NoError(t, err)

	amount, err := money.Zero(money.BRL)
	require.NoError(t, err)

	err = w.Credit(amount, now.Add(time.Hour))

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInvalidAmount)

	assert.Equal(t, "100.00", w.Balance().String())
	assert.Equal(t, int64(1), w.Version())
	assert.Equal(t, now, w.UpdatedAt())
}

func TestWalletCreditInvalidAmount(t *testing.T) {
	now := time.Now()
	balance, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)
	w, err := New("wallet-1", "player-1", money.BRL, balance, now)
	require.NoError(t, err)
	negativeAmount, err := money.NewFromMinorUnits(-100, money.BRL)
	require.NoError(t, err)

	tests := []struct {
		name   string
		amount money.Money
	}{
		{name: "negative amount", amount: negativeAmount},
		{name: "invalid amount", amount: money.Money{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := w.Credit(tt.amount, now.Add(time.Hour))

			require.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalidAmount)
			assert.Equal(t, "100.00", w.Balance().String())
			assert.Equal(t, int64(1), w.Version())
			assert.Equal(t, now, w.UpdatedAt())
		})
	}
}

func TestWalletCreditOverflow(t *testing.T) {
	now := time.Now()
	balance, err := money.NewFromMinorUnits(1, money.BRL)
	require.NoError(t, err)
	w, err := New("wallet-1", "player-1", money.BRL, balance, now)
	require.NoError(t, err)
	amount, err := money.NewFromMinorUnits(int64(^uint64(0)>>1), money.BRL)
	require.NoError(t, err)

	err = w.Credit(amount, now.Add(time.Hour))

	require.Error(t, err)
	assert.ErrorIs(t, err, money.ErrAmountOverflow)
	assert.Equal(t, "0.01", w.Balance().String())
	assert.Equal(t, int64(1), w.Version())
	assert.Equal(t, now, w.UpdatedAt())
}

func TestWalletCreditCurrencyMismatch(t *testing.T) {
	now := time.Now()

	balance, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "player-1", money.BRL, balance, now)
	require.NoError(t, err)

	amount, err := money.Parse("10.00", money.USD)
	require.NoError(t, err)

	err = w.Credit(amount, now.Add(time.Hour))

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrCurrencyMismatch)

	assert.Equal(t, "100.00", w.Balance().String())
	assert.Equal(t, int64(1), w.Version())
	assert.Equal(t, now, w.UpdatedAt())
}

func TestWalletDebit(t *testing.T) {
	t.Run("should debit wallet", func(t *testing.T) {
		createdAt := time.Date(2026, time.October, 4, 10, 0, 0, 0, time.UTC)
		debitedAt := createdAt.Add(time.Hour)

		balance, err := money.Parse("100.00", money.BRL)
		require.NoError(t, err)

		w, err := New("wallet-1", "player-1", money.BRL, balance, createdAt)
		require.NoError(t, err)

		amount, err := money.Parse("25.50", money.BRL)
		require.NoError(t, err)

		err = w.Debit(amount, debitedAt)

		require.NoError(t, err)

		assert.Equal(t, "74.50", w.Balance().String())
		assert.Equal(t, int64(2), w.Version())
		assert.Equal(t, debitedAt, w.UpdatedAt())
	})
}

func TestWalletDebitFullBalance(t *testing.T) {
	now := time.Now()

	balance, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "player-1", money.BRL, balance, now)
	require.NoError(t, err)

	amount, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	err = w.Debit(amount, now.Add(time.Hour))

	require.NoError(t, err)

	assert.True(t, w.Balance().IsZero())
	assert.Equal(t, int64(2), w.Version())
}

func TestWalletDebitInsufficientFunds(t *testing.T) {
	now := time.Now()

	balance, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "player-1", money.BRL, balance, now)
	require.NoError(t, err)

	amount, err := money.Parse("150.00", money.BRL)
	require.NoError(t, err)

	err = w.Debit(amount, now.Add(time.Hour))

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInsufficientFunds)

	// Nenhum estado pode ter sido alterado.
	assert.Equal(t, "100.00", w.Balance().String())
	assert.Equal(t, int64(1), w.Version())
	assert.Equal(t, now, w.UpdatedAt())
}

func TestWalletDebitInvalidAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
	}{
		{
			name:   "should reject zero",
			amount: 0,
		},
		{
			name:   "should reject negative amount",
			amount: -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()

			balance, err := money.Parse("100.00", money.BRL)
			require.NoError(t, err)

			w, err := New("wallet-1", "player-1", money.BRL, balance, now)
			require.NoError(t, err)

			amount, err := money.NewFromMinorUnits(tt.amount, money.BRL)
			require.NoError(t, err)

			err = w.Debit(amount, now.Add(time.Hour))

			require.Error(t, err)

			assert.ErrorIs(t, err, ErrInvalidAmount)
			assert.Equal(t, "100.00", w.Balance().String())
			assert.Equal(t, int64(1), w.Version())
			assert.Equal(t, now, w.UpdatedAt())
		})
	}
}

func TestWalletDebitCurrencyMismatch(t *testing.T) {
	now := time.Now()

	balance, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)

	w, err := New("wallet-1", "player-1", money.BRL, balance, now)
	require.NoError(t, err)

	amount, err := money.Parse("10.00", money.USD)
	require.NoError(t, err)

	err = w.Debit(amount, now.Add(time.Hour))

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrCurrencyMismatch)

	assert.Equal(t, "100.00", w.Balance().String())
	assert.Equal(t, int64(1), w.Version())
	assert.Equal(t, now, w.UpdatedAt())
}

func TestWalletOperationsRejectInvalidTimestamps(t *testing.T) {
	now := time.Date(2026, time.October, 4, 10, 0, 0, 0, time.UTC)
	balance, err := money.Parse("100.00", money.BRL)
	require.NoError(t, err)
	amount, err := money.Parse("10.00", money.BRL)
	require.NoError(t, err)

	for _, operation := range []struct {
		name string
		run  func(*Wallet, money.Money, time.Time) error
	}{
		{name: "credit", run: (*Wallet).Credit},
		{name: "debit", run: (*Wallet).Debit},
	} {
		t.Run(operation.name, func(t *testing.T) {
			for _, timestamp := range []struct {
				name string
				at   time.Time
			}{
				{name: "zero", at: time.Time{}},
				{name: "before last update", at: now.Add(-time.Nanosecond)},
			} {
				t.Run(timestamp.name, func(t *testing.T) {
					w, err := New("wallet-1", "player-1", money.BRL, balance, now)
					require.NoError(t, err)

					err = operation.run(w, amount, timestamp.at)

					require.Error(t, err)
					assert.ErrorIs(t, err, ErrInvalidTimestamp)
					assert.Equal(t, "100.00", w.Balance().String())
					assert.Equal(t, int64(1), w.Version())
					assert.Equal(t, now, w.UpdatedAt())
				})
			}
		})
	}
}
