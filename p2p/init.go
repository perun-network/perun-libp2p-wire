package p2p

import (
	"math/rand"

	"perun.network/go-perun/wire"
	"perun.network/go-perun/wire/test"
)

func init() {
	wire.SetNewAddressFunc(func() wire.Address {
		return NewAddress("")
	})
	// Setup for testing purposes.
	test.SetNewRandomAddress(func(rng *rand.Rand) wire.Address {
		return NewRandomAddress(rng)
	})
	// Setup for testing purposes.
	test.SetNewRandomAccount(func(rng *rand.Rand) wire.Account {
		return NewRandomAccount(rng)
	})
}
