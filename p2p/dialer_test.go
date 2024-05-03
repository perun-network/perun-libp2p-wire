package p2p

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"perun.network/go-perun/wire"
	perunio "perun.network/go-perun/wire/perunio/serializer"

	ctxtest "polycry.pt/poly-go/context/test"
	pkgtest "polycry.pt/poly-go/test"
)

func TestNewDialer(t *testing.T) {
	h := getHost(t)

	d := NewP2PDialer(h)
	assert.NotNil(t, d)
	d.Close()
}

func TestDialer_Register(t *testing.T) {
	rng := pkgtest.Prng(t)
	addr := NewRandomAddress(rng)
	key := wire.Key(addr)
	h := getHost(t)
	d := NewP2PDialer(h)
	defer d.Close()

	_, ok := d.get(key)
	require.False(t, ok)

	d.Register(addr, "p2pAddress")

	host, ok := d.get(key)
	assert.True(t, ok)
	assert.Equal(t, host, "p2pAddress")
}

func TestDialer_Dial(t *testing.T) {
	timeout := 1000 * time.Millisecond
	rng := pkgtest.Prng(t)

	lHost := getHost(t)
	lAddr := lHost.Address()
	lP2PAddr, err := getHostMA(lHost)
	require.NoError(t, err)
	listener := NewP2PListener(lHost)
	defer listener.Close()

	dHost := getHost(t)
	dAddr := dHost.Address()
	dialer := NewP2PDialer(dHost)
	dialer.Register(lAddr, lP2PAddr.String())
	defer dialer.Close()

	t.Run("happy", func(t *testing.T) {
		e := &wire.Envelope{
			Sender:    dAddr,
			Recipient: lAddr,
			Msg:       wire.NewPingMsg()}
		ct := pkgtest.NewConcurrent(t)

		go ct.Stage("accept", func(rt pkgtest.ConcT) {
			conn, err := listener.Accept()
			assert.NoError(t, err)
			require.NotNil(rt, conn)

			re, err := conn.Recv()
			assert.NoError(t, err)
			assert.Equal(t, re, e)
		})

		ct.Stage("dial", func(rt pkgtest.ConcT) {
			ctxtest.AssertTerminates(t, timeout, func() {
				conn, err := dialer.Dial(context.Background(), lAddr, perunio.Serializer())
				assert.NoError(t, err)
				require.NotNil(rt, conn)

				assert.NoError(t, conn.Send(e))
			})
		})

		ct.Wait("dial", "accept")
	})

	t.Run("aborted context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ctxtest.AssertTerminates(t, timeout, func() {
			conn, err := dialer.Dial(ctx, lAddr, perunio.Serializer())
			assert.Nil(t, conn)
			assert.Error(t, err)
		})
	})

	t.Run("unknown host", func(t *testing.T) {
		noHostAddr := NewRandomAddress(rng)
		dialer.Register(noHostAddr, "no such host")

		ctxtest.AssertTerminates(t, timeout, func() {
			conn, err := dialer.Dial(context.Background(), noHostAddr, perunio.Serializer())
			assert.Nil(t, conn)
			assert.Error(t, err)
		})
	})

	t.Run("unknown address", func(t *testing.T) {
		ctxtest.AssertTerminates(t, timeout, func() {
			unknownAddr := NewRandomAddress(rng)
			conn, err := dialer.Dial(context.Background(), unknownAddr, perunio.Serializer())
			assert.Error(t, err)
			assert.Nil(t, conn)
		})
	})
}

func TestDialer_Close(t *testing.T) {
	t.Run("double close", func(t *testing.T) {
		h := getHost(t)
		d := NewP2PDialer(h)

		assert.NoError(t, d.Close(), "first close must not return error")
		assert.Error(t, d.Close(), "second close must result in error")
	})
}
