package wager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateIsValid(t *testing.T) {
	tests := []struct {
		state State
		valid bool
	}{
		{StatePending, true},
		{StatePendingReference, true},
		{StateProcessed, true},
		{StateRejected, true},
		{StateFailed, true},
		{State("INVALID"), false},
		{State(""), false},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.valid, tt.state.IsValid())
	}
}

func TestStateIsTerminal(t *testing.T) {
	tests := []struct {
		state    State
		terminal bool
	}{
		{StatePending, false},
		{StatePendingReference, false},
		{StateProcessed, true},
		{StateRejected, true},
		{StateFailed, true},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.terminal, tt.state.IsTerminal())
	}
}
