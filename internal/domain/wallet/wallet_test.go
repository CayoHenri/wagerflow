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
