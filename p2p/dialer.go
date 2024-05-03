package p2p

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"perun.network/go-perun/wire"
	wirenet "perun.network/go-perun/wire/net"
	pkgsync "polycry.pt/poly-go/sync"
)

// Dialer is a dialer for p2p connections.
type Dialer struct {
	mutex  sync.RWMutex // Protects peers.
	peers  map[wire.AddrKey]string
	host   host.Host
	closer pkgsync.Closer
}

// NewP2PDialer creates a new dialer for the given account.
func NewP2PDialer(acc *Account) *Dialer {
	return &Dialer{
		host:  acc,
		peers: make(map[wire.AddrKey]string),
	}
}

// Dial implements Dialer.Dial().
func (d *Dialer) Dial(ctx context.Context, addr wire.Address, serializer wire.EnvelopeSerializer) (wirenet.Conn, error) {
	peerMA, ok := d.get(wire.Key(addr))
	if !ok {
		return nil, errors.New("peer not found")
	}

	_peerMA, err := ma.NewMultiaddr(peerMA)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse multiaddress of peer")
	}

	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(_peerMA)
	if err != nil {
		return nil, errors.Wrap(err, "converting peer multiaddress to address info")
	}
	if err := d.host.Connect(ctx, *peerAddrInfo); err != nil {
		return nil, errors.Wrap(err, "failed to dial peer: failed to connecting to peer")
	}

	s, err := d.host.NewStream(ctx, peerAddrInfo.ID, "/client")
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial peer: failed to creating a new stream")
	}

	return wirenet.NewIoConn(s, serializer), nil
}

// Register registers a p2p multiaddress for a peer wire address.
func (d *Dialer) Register(addr wire.Address, p2pAddress string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.peers[wire.Key(addr)] = p2pAddress
}

// Close closes the dialer by closing the underlying libp2p host.
func (d *Dialer) Close() error {
	if err := d.closer.Close(); err != nil {
		return err
	}
	return d.host.Close()
}

// get returns the p2p multiaddress for the given address if registered.
func (d *Dialer) get(addr wire.AddrKey) (p2pAddress string, ok bool) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	p2pAddress, ok = d.peers[addr]
	return
}
