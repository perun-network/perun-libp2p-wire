package p2p

import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/pkg/errors"
	"perun.network/go-perun/wire"
	wirenet "perun.network/go-perun/wire/net"
	pkgsync "polycry.pt/poly-go/sync"
)

// Listener is a listener for p2p connections.
type Listener struct {
	host    host.Host
	streams chan network.Stream
	done    chan struct{}
	closer  pkgsync.Closer
}

// NewP2PListener creates a new listener for the given account.
func NewP2PListener(acc *Account) *Listener {
	listener := Listener{
		host:    acc,
		streams: make(chan network.Stream),
		done:    make(chan struct{}),
	}
	acc.SetStreamHandler("/client", func(s network.Stream) {
		listener.streams <- s
	})
	return &listener
}

// Accept implements Listener.Accept().
func (l *Listener) Accept(serializer wire.EnvelopeSerializer) (wirenet.Conn, error) {
	select {
	case s := <-l.streams:
		return wirenet.NewIoConn(s, serializer), nil
	case <-l.done:
		return nil, errors.New("listener is closed")
	}
}

// Close closes the Listener by closing the done channel.
func (l *Listener) Close() error {
	if err := l.closer.Close(); err != nil {
		return err
	}

	close(l.done)
	return nil
}
