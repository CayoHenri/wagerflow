package wager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKindIsValid(t *testing.T) {
	tests := []struct {
		name  string
		kind  Kind
		valid bool
	}{
		{"opening", KindOpening, true},
		{"bet", KindBet, true},
		{"win", KindWin, true},
		{"loss", KindLoss, true},
		{"refund", KindRefund, true},
		{"rollback", KindRollback, true},
		{"invalid", Kind("INVALID"), false},
		{"empty", Kind(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.kind.IsValid())
		})
	}
}
