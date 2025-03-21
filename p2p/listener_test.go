package p2p

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	perunio "perun.network/go-perun/wire/perunio/serializer"

	ctxtest "polycry.pt/poly-go/context/test"
	pkgtest "polycry.pt/poly-go/test"
)

func TestNewListener(t *testing.T) {
	rng := pkgtest.Prng(t)
	h := getHost(rng)
	l := NewP2PListener(h)
	defer l.Close()
	assert.NotNil(t, l)
}

func TestListener_Close(t *testing.T) {
	t.Run("double close", func(t *testing.T) {
		rng := pkgtest.Prng(t)
		h := getHost(rng)
		l := NewP2PListener(h)
		assert.NoError(t, l.Close(), "first close must not return error")
		assert.Error(t, l.Close(), "second close must result in error")
	})
}

func TestListener_Accept(t *testing.T) {
	// Happy case already tested in TestDialer_Dial.
	rng := pkgtest.Prng(t)
	h := getHost(rng)
	timeout := 100 * time.Millisecond
	t.Run("timeout", func(t *testing.T) {
		l := NewP2PListener(h)
		defer l.Close()

		ctxtest.AssertNotTerminates(t, timeout, func() {
			l.Accept(perunio.Serializer())
		})
	})

	t.Run("closed", func(t *testing.T) {
		l := NewP2PListener(h)
		l.Close()

		ctxtest.AssertTerminates(t, timeout, func() {
			conn, err := l.Accept(perunio.Serializer())
			assert.Nil(t, conn)
			assert.Error(t, err)
		})
	})
}
