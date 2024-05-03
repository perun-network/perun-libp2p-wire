package p2p_test

import (
	"math/rand"
	"testing"

	"github.com/perun-network/perun-libp2p-wire/p2p"
	"perun.network/go-perun/wire"
	"perun.network/go-perun/wire/test"
)

func TestAddress(t *testing.T) {
	test.TestAddressImplementation(t, func() wire.Address {
		return p2p.NewAddress("")
	}, func(rng *rand.Rand) wire.Address {
		return p2p.NewRandomAddress(rng)
	})
}
