package p2p

import (
	"context"
	"crypto/sha256"
	"math/rand"
	"net"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"perun.network/go-perun/wire"
)

const relayID = "QmVCPfUMr98PaaM8qbAQBgJ9jqc7XHpGp7AsyragdFDmgm"

// Account represents a libp2p wire account.
type Account struct {
	host.Host
	relayAddr  string
	privateKey crypto.PrivKey
}

// Address returns the account's address.
func (acc *Account) Address() wire.Address {
	return &Address{acc.ID()}
}

// Sign signs the given message with the account's private key.
func (acc *Account) Sign(data []byte) ([]byte, error) {
	// Extract the private key from the account.
	if acc.privateKey == nil {
		return nil, errors.New("private key not set")
	}
	hashed := sha256.Sum256(data)

	signature, err := acc.privateKey.Sign(hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil

}

// NewRandomAccount generates a new random account.
func NewRandomAccount(rng *rand.Rand) *Account {
	id, err := peer.Decode(relayID)
	if err != nil {
		err = errors.WithMessage(err, "decoding peer id of relay server")
		return nil
	}

	// Get the IP address of the relay server.
	ip, err := net.LookupIP("relay.perun.network")
	if err != nil {
		panic(errors.WithMessage(err, "looking up IP address of relay.perun.network"))
	}
	relayAddr := "/ip4/" + ip[0].String() + "/tcp/5574"

	relayMultiaddr, err := ma.NewMultiaddr(relayAddr)
	if err != nil {
		panic(errors.WithMessage(err, "parsing relay multiadress"))
	}

	relayInfo := peer.AddrInfo{
		ID:    id,
		Addrs: []ma.Multiaddr{relayMultiaddr},
	}

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rng)
	if err != nil {
		panic(err)
	}

	// Construct a new libp2p client for our relay-server.
	// Identity(prvKey)		- Use a RSA private key to generate the ID of the host.
	// EnableRelay()		- Enable relay system and configures itself as a node behind a NAT
	client, err := libp2p.New(
		context.Background(),
		libp2p.Identity(prvKey),
		libp2p.EnableRelay(),
	)
	if err != nil {
		panic(err)
	}

	if err := client.Connect(context.Background(), relayInfo); err != nil {
		panic(errors.WithMessage(err, "connecting to the relay server"))
	}
	return &Account{client, relayAddr, prvKey}
}
