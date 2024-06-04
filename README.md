# Perun LibP2P Wire

## Introduction
This repository serves as an implementation of the `wire` package in [Go-Perun] v0.10.6, which use the technology of [Go-Libp2p] (v0.13.0). This implementation realise the `Account` as libp2p's `host.Host` and `Address` as libp2p's `PeerID`.

The Account will automatically connect to [Perun-Relay-Server] after creation, allowing for peer-to-peer connection with other clients, if their Addresses (PeerIDs) are already known.

## Usage
Perun LibP2P Wire could be used to create bus connection for off-chain communication for every Perun Client. 

Example usage:
```go
    import "github.com/perun-network/perun-libp2p-wire/p2p"

	// Create bus and listener.
	wireAcc := p2p.NewRandomAccount(rand.New(rand.NewSource(time.Now().UnixNano())))

	net, err := p2p.NewP2PBus(wireAcc)
	if err != nil {
		panic(errors.Wrap(err, "creating p2p net"))
	}

    bus := net.Bus
	listener := net.Listener

	// Create and start client.
	c, err := perunclient.New(
		wireAcc.Address(),
        bus,
        ... // Other fields
	)
	if err != nil {
		panic(err)
	}
	go bus.Listen(listener)
```

## Constraint
The `Dialer` requires the other peers to be already "registered" to connect with them. Before dialing, the `Register` must be called. 

Example (register peer before propose a channel with them):
````go
// OpenChannel opens a new channel with the specified peer and funding.
func (c *PaymentClient) OpenChannel(peer wire.Address, peerID string, amount float64) *PaymentChannel {
	....
	```
	c.net.Dialer.Register(peer, peerID)
	```
	... 

	// Prepare the channel proposal by defining the channel parameters.
	challengeDuration := uint64(10) // On-chain challenge duration in seconds.
	proposal, err := client.NewLedgerChannelProposal(
		challengeDuration,
		c.account,
		initAlloc,
		participants,
	)
	if err != nil {
		panic(err)
	}

	// Send the proposal.
	ch, err := c.perunClient.ProposeChannel(context.TODO(), proposal)
	if err != nil {
		panic(err)
	}

	// Start the on-chain event watcher. It automatically handles disputes.
	c.startWatching(ch)

	return newPaymentChannel(ch, c.currency)
}
````

## Address Ressolver
A default address resolver was already built-in on the Perun-Relay-Server. You can use the provided APIs in order to `Register`, `Query`, `Deregister` your On-chain (L1) Address (Implementation of [Go-Perun] `wallet.Address`) to get the peer's `wire.Address` (`Peer.ID` of [Go-Libp2p])

**Example:**
````go
// Should be used in the initialization of Perun-Client.
err := acc.RegisterOnChainAddress(onChainAddr)


// Query the peer's wire address, given its on-chain address.
peerID, err := acc.QueryOnChainAddress(peerOnChainAddr)


// Deregister the on-chain address, to be used before closing Perun-Client,
err = acc.DeregisterOnChainAddress(onChainAddr)

````

## Test
Some unit tests are provided:
```
go test -v ./...
```


[Go-Perun]: https://github.com/hyperledger-labs/go-perun
[Go-Libp2p]: https://pkg.go.dev/github.com/libp2p/go-libp2p@v0.13.0
[Perun-Relay-Server]: https://github.com/perun-network/perun-relay