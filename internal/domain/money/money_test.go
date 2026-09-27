package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should create money with amount and currency", func(t *testing.T) {
		m := New(1025, BRL)

		assert.Equal(t, int64(1025), m.Amount())
		assert.Equal(t, BRL, m.Currency())
	})
}
