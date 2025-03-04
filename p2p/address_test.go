package p2p_test

import (
	"math/rand"
	"testing"

	"github.com/perun-network/perun-libp2p-wire/p2p"
	"github.com/stretchr/testify/assert"
	"perun.network/go-perun/wire"
	"perun.network/go-perun/wire/test"
	pkgtest "polycry.pt/poly-go/test"
)

func TestAddress(t *testing.T) {
	test.TestAddressImplementation(t, func() wire.Address {
		return p2p.NewAddress("")
	}, func(rng *rand.Rand) wire.Address {
		return p2p.NewRandomAddress(rng)
	})
}

func TestSignature(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := p2p.NewRandomAccount(rng)
	assert.NotNil(t, acc)
	defer acc.Close()

	msg := []byte("test message")
	sig, err := acc.Sign(msg)
	assert.NoError(t, err)

	addr := acc.Address()
	assert.NoError(t, addr.Verify(msg, sig))
}
