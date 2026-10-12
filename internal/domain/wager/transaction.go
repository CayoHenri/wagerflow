package wager

import (
	"strings"
	"time"

	"github.com/CayoHenri/wagerflow/internal/domain/money"
)

type WagerTransaction struct {
	id                             string
	externalTransactionID          string
	providerID                     string
	idempotencyKey                 string
	payloadHash                    string
	walletID                       string
	playerID                       string
	roundID                        string
	gameID                         string
	kind                           Kind
	state                          State
	amount                         money.Money
	referenceExternalTransactionID string
	resolvedReferenceID            string
	failureCode                    string
	createdAt                      time.Time
	updatedAt                      time.Time
}

func New(
	id string,
	externalTransactionID string,
	providerID string,
	idempotencyKey string,
	payloadHash string,
	walletID string,
	playerID string,
	roundID string,
	gameID string,
	kind Kind,
	amount money.Money,
	referenceExternalTransactionID string,
	now time.Time,
) (*WagerTransaction, error) {
	id = strings.TrimSpace(id)
	externalTransactionID = strings.TrimSpace(externalTransactionID)
	providerID = strings.TrimSpace(providerID)
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	payloadHash = strings.TrimSpace(payloadHash)
	walletID = strings.TrimSpace(walletID)
	playerID = strings.TrimSpace(playerID)
	roundID = strings.TrimSpace(roundID)
	gameID = strings.TrimSpace(gameID)
	referenceExternalTransactionID = strings.TrimSpace(referenceExternalTransactionID)

	if id == "" {
		return nil, ErrInvalidID
	}

	if externalTransactionID == "" {
		return nil, ErrInvalidExternalTransactionID
	}

	if providerID == "" {
		return nil, ErrInvalidProviderID
	}

	if idempotencyKey == "" {
		return nil, ErrInvalidIdempotencyKey
	}

	if payloadHash == "" {
		return nil, ErrInvalidPayloadHash
	}

	if walletID == "" {
		return nil, ErrInvalidWalletID
	}

	if playerID == "" {
		return nil, ErrInvalidPlayerID
	}

	if roundID == "" {
		return nil, ErrInvalidRoundID
	}

	if gameID == "" {
		return nil, ErrInvalidGameID
	}

	if !kind.IsValid() {
		return nil, ErrInvalidKind
	}

	if !amount.IsValid() {
		return nil, ErrInvalidAmount
	}

	if now.IsZero() {
		return nil, ErrInvalidTimestamp
	}

	return &WagerTransaction{
		id:                             id,
		externalTransactionID:          externalTransactionID,
		providerID:                     providerID,
		idempotencyKey:                 idempotencyKey,
		payloadHash:                    payloadHash,
		walletID:                       walletID,
		playerID:                       playerID,
		roundID:                        roundID,
		gameID:                         gameID,
		kind:                           kind,
		state:                          StatePending,
		amount:                         amount,
		referenceExternalTransactionID: referenceExternalTransactionID,
		createdAt:                      now,
		updatedAt:                      now,
	}, nil
}

func (t WagerTransaction) ID() string {
	return t.id
}

func (t WagerTransaction) ExternalTransactionID() string {
	return t.externalTransactionID
}

func (t WagerTransaction) ProviderID() string {
	return t.providerID
}

func (t WagerTransaction) IdempotencyKey() string {
	return t.idempotencyKey
}

func (t WagerTransaction) PayloadHash() string {
	return t.payloadHash
}

func (t WagerTransaction) WalletID() string {
	return t.walletID
}

func (t WagerTransaction) PlayerID() string {
	return t.playerID
}

func (t WagerTransaction) RoundID() string {
	return t.roundID
}

func (t WagerTransaction) GameID() string {
	return t.gameID
}

func (t WagerTransaction) Kind() Kind {
	return t.kind
}

func (t WagerTransaction) State() State {
	return t.state
}

func (t WagerTransaction) Amount() money.Money {
	return t.amount
}

func (t WagerTransaction) ReferenceExternalTransactionID() string {
	return t.referenceExternalTransactionID
}

func (t WagerTransaction) ResolvedReferenceID() string {
	return t.resolvedReferenceID
}

func (t WagerTransaction) FailureCode() string {
	return t.failureCode
}

func (t WagerTransaction) CreatedAt() time.Time {
	return t.createdAt
}

func (t WagerTransaction) UpdatedAt() time.Time {
	return t.updatedAt
}
