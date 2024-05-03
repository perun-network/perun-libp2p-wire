package p2p

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ctxtest "polycry.pt/poly-go/context/test"
)

func TestNewListener(t *testing.T) {
	h := getHost(t)
	l := NewP2PListener(h)
	defer l.Close()
	assert.NotNil(t, l)
}

func TestListener_Close(t *testing.T) {
	t.Run("double close", func(t *testing.T) {
		h := getHost(t)
		l := NewP2PListener(h)
		assert.NoError(t, l.Close(), "first close must not return error")
		assert.Error(t, l.Close(), "second close must result in error")
	})
}

func TestListener_Accept(t *testing.T) {
	// Happy case already tested in TestDialer_Dial.

	h := getHost(t)
	timeout := 100 * time.Millisecond
	t.Run("timeout", func(t *testing.T) {
		l := NewP2PListener(h)
		defer l.Close()

		ctxtest.AssertNotTerminates(t, timeout, func() {
			l.Accept()
		})
	})

	t.Run("closed", func(t *testing.T) {
		l := NewP2PListener(h)
		l.Close()

		ctxtest.AssertTerminates(t, timeout, func() {
			conn, err := l.Accept()
			assert.Nil(t, conn)
			assert.Error(t, err)
		})
	})
}
