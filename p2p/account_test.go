package p2p

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sim_wallet "perun.network/go-perun/backend/sim/wallet"
	"perun.network/go-perun/wallet"
	pkgtest "polycry.pt/poly-go/test"
)

func TestNewAccount(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
	defer acc.Close()
}

func getHost(t *testing.T) *Account {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
	return acc
}

func TestAddressBookRegister(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
	defer acc.Close()

	onChainAddr := sim_wallet.NewRandomAddress(rng)

	err := acc.RegisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
}

func TestAddressBookRegisterEmptyAddress(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)

	defer acc.Close()

	emptyAddr := &sim_wallet.Address{}

	assert.Panics(t, func() { acc.RegisterOnChainAddress(emptyAddr) })

	var nilAddr wallet.Address
	err := acc.RegisterOnChainAddress(nilAddr)
	assert.Error(t, err)
}

func TestAddressBookDeregister(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
	defer acc.Close()

	onChainAddr := sim_wallet.NewRandomAddress(rng)

	err := acc.RegisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	err = acc.DeregisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	// Trying to query it again will fail
	_, err = acc.QueryOnChainAddress(onChainAddr)
	assert.Error(t, err)
}

func TestAddressBookDeregisterPeer(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
	defer acc.Close()

	peer := NewRandomAccount(rng)
	assert.NotNil(t, peer)
	defer peer.Close()

	onChainAddr := sim_wallet.NewRandomAddress(rng)
	peerOnChainAddr := sim_wallet.NewRandomAddress(rng)

	err := acc.RegisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	time.Sleep(1 * time.Millisecond)

	err = peer.RegisterOnChainAddress(peerOnChainAddr)
	assert.NoError(t, err)

	err = acc.DeregisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	// Trying to deregister the peer's address will not fail, but the server will not allow it.
	err = acc.DeregisterOnChainAddress(peerOnChainAddr)
	assert.NoError(t, err)

	// Trying to query it again will be okay
	peerID, err := acc.QueryOnChainAddress(peerOnChainAddr)
	assert.NoError(t, err)

	addr := peer.Address()
	assert.Equal(t, peerID, addr)

	err = peer.DeregisterOnChainAddress(peerOnChainAddr)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
}

func TestAddressBookQuery_Fail(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
	defer acc.Close()

	onChainAddr := sim_wallet.NewRandomAddress(rng)

	_, err := acc.QueryOnChainAddress(onChainAddr)
	assert.Error(t, err)
}

func TestAddressBookQuery(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
	defer acc.Close()

	onChainAddr := sim_wallet.NewRandomAddress(rng)

	err := acc.RegisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	time.Sleep(10 * time.Millisecond)
	peerID, err := acc.QueryOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	addr := acc.Address()
	assert.Equal(t, peerID, addr)

	err = acc.DeregisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
}

func TestAddressBookQueryPeer(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
	defer acc.Close()

	peer := NewRandomAccount(rng)
	assert.NotNil(t, peer)
	defer peer.Close()

	onChainAddr := sim_wallet.NewRandomAddress(rng)
	peerOnChainAddr := sim_wallet.NewRandomAddress(rng)

	err := acc.RegisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	err = peer.RegisterOnChainAddress(peerOnChainAddr)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
	peerID, err := acc.QueryOnChainAddress(peerOnChainAddr)
	assert.NoError(t, err)

	addr := peer.Address()
	assert.Equal(t, peerID, addr)

	err = acc.DeregisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	err = acc.DeregisterOnChainAddress(peerOnChainAddr)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
}

func TestAddressBookRegisterQueryMultiple(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)
	defer acc.Close()

	onChainAddr := sim_wallet.NewRandomAddress(rng)
	onChainAddr2 := sim_wallet.NewRandomAddress(rng)

	err := acc.RegisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	err = acc.RegisterOnChainAddress(onChainAddr2)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)

	accID, err := acc.QueryOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	accID2, err := acc.QueryOnChainAddress(onChainAddr2)
	assert.NoError(t, err)

	addr := acc.Address()
	assert.Equal(t, accID, addr)
	assert.Equal(t, accID2, addr)

	// Clean up
	err = acc.DeregisterOnChainAddress(onChainAddr)
	assert.NoError(t, err)

	err = acc.DeregisterOnChainAddress(onChainAddr2)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
}

// Test NewAccountFromPrivateKey
func TestNewAccountFromPrivateKey(t *testing.T) {
	rng := pkgtest.Prng(t)
	acc := NewRandomAccount(rng)
	assert.NotNil(t, acc)

	defer acc.Close()

	keyBytes, err := acc.MarshalPrivateKey()
	assert.NoError(t, err)

	acc2, err := NewAccountFromPrivateKeyBytes(keyBytes)
	assert.NoError(t, err)

	defer acc2.Close()

	assert.NotNil(t, acc2)
	assert.Equal(t, acc.ID(), acc2.ID())
	assert.Equal(t, acc.Address(), acc2.Address())
}
