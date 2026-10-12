package wager

import (
	"testing"
	"time"

	"github.com/CayoHenri/wagerflow/internal/domain/money"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWagerTransaction(t *testing.T) {
	now := time.Date(2026, time.October, 6, 13, 0, 0, 0, time.UTC)

	amount, err := money.Parse("25.00", money.BRL)
	require.NoError(t, err)

	transaction, err := New(
		"transaction-1",
		"external-1",
		"provider-1",
		"idempotency-1",
		"payload-hash-1",
		"wallet-1",
		"player-1",
		"round-1",
		"game-1",
		KindBet,
		amount,
		"",
		now,
	)

	require.NoError(t, err)
	require.NotNil(t, transaction)

	assert.Equal(t, "transaction-1", transaction.ID())
	assert.Equal(t, "external-1", transaction.ExternalTransactionID())
	assert.Equal(t, "provider-1", transaction.ProviderID())
	assert.Equal(t, "idempotency-1", transaction.IdempotencyKey())
	assert.Equal(t, "payload-hash-1", transaction.PayloadHash())
	assert.Equal(t, "wallet-1", transaction.WalletID())
	assert.Equal(t, "player-1", transaction.PlayerID())
	assert.Equal(t, "round-1", transaction.RoundID())
	assert.Equal(t, "game-1", transaction.GameID())

	assert.Equal(t, KindBet, transaction.Kind())
	assert.Equal(t, StatePending, transaction.State())

	assert.True(t, amount.Equal(transaction.Amount()))

	assert.Empty(t, transaction.ReferenceExternalTransactionID())
	assert.Empty(t, transaction.ResolvedReferenceID())
	assert.Empty(t, transaction.FailureCode())

	assert.Equal(t, now, transaction.CreatedAt())
	assert.Equal(t, now, transaction.UpdatedAt())
}

func TestNewRejectsInvalidIdentifiers(t *testing.T) {
	amount, err := money.Parse("10.00", money.BRL)
	require.NoError(t, err)

	tests := []struct {
		name                  string
		id                    string
		externalTransactionID string
		providerID            string
		idempotencyKey        string
		payloadHash           string
		walletID              string
		playerID              string
		roundID               string
		gameID                string
		expectedErr           error
	}{
		{
			name:                  "empty id",
			id:                    "",
			externalTransactionID: "external-1",
			providerID:            "provider-1",
			idempotencyKey:        "key-1",
			payloadHash:           "hash-1",
			walletID:              "wallet-1",
			playerID:              "player-1",
			roundID:               "round-1",
			gameID:                "game-1",
			expectedErr:           ErrInvalidID,
		},
		{
			name:                  "empty provider",
			id:                    "transaction-1",
			externalTransactionID: "external-1",
			providerID:            "",
			idempotencyKey:        "key-1",
			payloadHash:           "hash-1",
			walletID:              "wallet-1",
			playerID:              "player-1",
			roundID:               "round-1",
			gameID:                "game-1",
			expectedErr:           ErrInvalidProviderID,
		},
		{
			name:                  "empty external transaction id",
			id:                    "transaction-1",
			externalTransactionID: "",
			providerID:            "provider-1",
			idempotencyKey:        "key-1",
			payloadHash:           "hash-1",
			walletID:              "wallet-1",
			playerID:              "player-1",
			roundID:               "round-1",
			gameID:                "game-1",
			expectedErr:           ErrInvalidExternalTransactionID,
		},
		{
			name:                  "empty idempotency key",
			id:                    "transaction-1",
			externalTransactionID: "external-1",
			providerID:            "provider-1",
			idempotencyKey:        "",
			payloadHash:           "hash-1",
			walletID:              "wallet-1",
			playerID:              "player-1",
			roundID:               "round-1",
			gameID:                "game-1",
			expectedErr:           ErrInvalidIdempotencyKey,
		},
		{
			name:                  "empty payload hash",
			id:                    "transaction-1",
			externalTransactionID: "external-1",
			providerID:            "provider-1",
			idempotencyKey:        "key-1",
			payloadHash:           "",
			walletID:              "wallet-1",
			playerID:              "player-1",
			roundID:               "round-1",
			gameID:                "game-1",
			expectedErr:           ErrInvalidPayloadHash,
		},
		{
			name:                  "empty wallet id",
			id:                    "transaction-1",
			externalTransactionID: "external-1",
			providerID:            "provider-1",
			idempotencyKey:        "key-1",
			payloadHash:           "hash-1",
			walletID:              "",
			playerID:              "player-1",
			roundID:               "round-1",
			gameID:                "game-1",
			expectedErr:           ErrInvalidWalletID,
		},
		{
			name:                  "empty player id",
			id:                    "transaction-1",
			externalTransactionID: "external-1",
			providerID:            "provider-1",
			idempotencyKey:        "key-1",
			payloadHash:           "hash-1",
			walletID:              "wallet-1",
			playerID:              "",
			roundID:               "round-1",
			gameID:                "game-1",
			expectedErr:           ErrInvalidPlayerID,
		},
		{
			name:                  "empty round id",
			id:                    "transaction-1",
			externalTransactionID: "external-1",
			providerID:            "provider-1",
			idempotencyKey:        "key-1",
			payloadHash:           "hash-1",
			walletID:              "wallet-1",
			playerID:              "player-1",
			roundID:               "",
			gameID:                "game-1",
			expectedErr:           ErrInvalidRoundID,
		},
		{
			name:                  "empty game id",
			id:                    "transaction-1",
			externalTransactionID: "external-1",
			providerID:            "provider-1",
			idempotencyKey:        "key-1",
			payloadHash:           "hash-1",
			walletID:              "wallet-1",
			playerID:              "player-1",
			roundID:               "round-1",
			gameID:                "",
			expectedErr:           ErrInvalidGameID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := New(
				tt.id,
				tt.externalTransactionID,
				tt.providerID,
				tt.idempotencyKey,
				tt.payloadHash,
				tt.walletID,
				tt.playerID,
				tt.roundID,
				tt.gameID,
				KindBet,
				amount,
				"",
				time.Now().UTC(),
			)

			require.Error(t, err)

			assert.ErrorIs(t, err, tt.expectedErr)
			assert.Nil(t, transaction)
		})
	}
}

func TestNewRejectsInvalidKind(t *testing.T) {
	amount, err := money.Parse("10.00", money.BRL)
	require.NoError(t, err)

	transaction, err := New(
		"transaction-1",
		"external-1",
		"provider-1",
		"key-1",
		"hash-1",
		"wallet-1",
		"player-1",
		"round-1",
		"game-1",
		Kind("INVALID"),
		amount,
		"",
		time.Now().UTC(),
	)

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInvalidKind)
	assert.Nil(t, transaction)
}

func TestNewRejectsInvalidAmount(t *testing.T) {
	transaction, err := New(
		"transaction-1",
		"external-1",
		"provider-1",
		"key-1",
		"hash-1",
		"wallet-1",
		"player-1",
		"round-1",
		"game-1",
		KindBet,
		money.Money{},
		"",
		time.Now().UTC(),
	)

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInvalidAmount)
	assert.Nil(t, transaction)
}

func TestNewRejectsInvalidTimestamp(t *testing.T) {
	amount, err := money.Parse("10.00", money.BRL)
	require.NoError(t, err)

	transaction, err := New(
		"transaction-1",
		"external-1",
		"provider-1",
		"key-1",
		"hash-1",
		"wallet-1",
		"player-1",
		"round-1",
		"game-1",
		KindBet,
		amount,
		"",
		time.Time{},
	)

	require.Error(t, err)

	assert.ErrorIs(t, err, ErrInvalidTimestamp)
	assert.Nil(t, transaction)
}
