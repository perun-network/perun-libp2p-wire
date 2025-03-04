package p2p

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	swarm "github.com/libp2p/go-libp2p/p2p/net/swarm"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"perun.network/go-perun/wire"
	wirenet "perun.network/go-perun/wire/net"
	pkgsync "polycry.pt/poly-go/sync"
)

// Dialer is a dialer for p2p connections.
type Dialer struct {
	mutex     sync.RWMutex // Protects peers.
	peers     map[wire.AddrKey]string
	host      host.Host
	relayID   string
	relayAddr string
	closer    pkgsync.Closer
}

// NewP2PDialer creates a new dialer for the given account.
func NewP2PDialer(acc *Account, relayID string, relayAddr string) *Dialer {
	return &Dialer{
		host:      acc,
		relayID:   relayID,
		relayAddr: relayAddr,
		peers:     make(map[wire.AddrKey]string),
	}
}

// Dial implements Dialer.Dial().
func (d *Dialer) Dial(ctx context.Context, addr wire.Address, serializer wire.EnvelopeSerializer) (wirenet.Conn, error) {
	peerID, ok := d.get(wire.Key(addr))
	if !ok {
		return nil, errors.New("failed to dial peer: peer ID not found")
	}

	_peerID, err := peer.Decode(peerID)
	if err != nil {
		return nil, errors.Wrap(err, "peer ID is not valid")
	}

	d.host.Network().(*swarm.Swarm).Backoff().Clear(_peerID)
	fullAddr := d.relayAddr + "/p2p/" + d.relayID + "/p2p-circuit/p2p/" + _peerID.String()
	peerMultiAddr, err := ma.NewMultiaddr(fullAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse multiaddress of peer")
	}

	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMultiAddr)
	if err != nil {
		return nil, errors.Wrap(err, "converting peer multiaddress to address info")
	}
	if err := d.host.Connect(ctx, *peerAddrInfo); err != nil {
		return nil, errors.Wrap(err, "failed to dial peer: failed to connecting to peer")
	}

	s, err := d.host.NewStream(network.WithAllowLimitedConn(ctx, "client"), peerAddrInfo.ID, "/client")
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial peer: failed to creating a new stream")
	}

	return wirenet.NewIoConn(s, serializer), nil
}

// Register registers a p2p peer id for a peer wire address.
func (d *Dialer) Register(addr wire.Address, peerID string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.peers[wire.Key(addr)] = peerID
}

// Close closes the dialer by closing the underlying libp2p host.
func (d *Dialer) Close() error {
	if err := d.closer.Close(); err != nil {
		return err
	}
	return d.host.Close()
}

// get returns the p2p multiaddress for the given address if registered.
func (d *Dialer) get(addr wire.AddrKey) (peerID string, ok bool) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	peerID, ok = d.peers[addr]
	return
}
