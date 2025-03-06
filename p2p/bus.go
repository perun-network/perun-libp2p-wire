package p2p

import (
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire"
	"perun.network/go-perun/wire/net"
	perunio "perun.network/go-perun/wire/perunio/serializer"
)

// Net contains the client's components for the P2P communication.
type Net struct {
	*net.Bus
	*Listener
	*Dialer
	PeerID string
}

// NewP2PBus creates a dialer, listener, and a bus for the given account `acc`
// and includes them in the returned P2P Net.
func NewP2PBus(backendID wallet.BackendID, acc *Account) (*Net, error) {
	listener := NewP2PListener(acc)
	dialer := NewP2PDialer(acc, relayID)

	id := make(map[wallet.BackendID]wire.Account)
	id[backendID] = acc

	bus := net.NewBus(id, dialer, perunio.Serializer())
	return &Net{Bus: bus, Dialer: dialer, Listener: listener, PeerID: acc.relayAddr}, nil
}
