package p2p

import (
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"perun.network/go-perun/wire/net"
	perunio "perun.network/go-perun/wire/perunio/serializer"
)

// Net contains the client's components for the P2P communication.
type Net struct {
	*net.Bus
	*Listener
	*Dialer
	MultiAddr string
}

// NewP2PBus creates a dialer, listener, and a bus for the given account `acc`
// and includes them in the returned P2P Net.
func NewP2PBus(acc *Account) (*Net, error) {
	listener := NewP2PListener(acc)
	dialer := NewP2PDialer(acc)
	bus := net.NewBus(acc, dialer, perunio.Serializer())

	multiAddr, err := getHostMA(acc)
	if err != nil {
		return nil, err
	}

	return &Net{Bus: bus, Dialer: dialer, Listener: listener, MultiAddr: multiAddr.String()}, nil
}

// getHostMA returns the first multiaddress of the given host.
func getHostMA(host host.Host) (multiaddr.Multiaddr, error) {
	peerInfo := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	return addrs[0], errors.Wrap(err, "converting peer info to multiaddress")
}
