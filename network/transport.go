package network

import "net"

// NetAddress represents a network address string.
type NetAddress string

// Transport defines methods for network communication.
type Transport interface {
	// Consume returns a channel to consume incoming RPC messages.
	Consume() <-chan RPC

	// Connect establishes a connection to another Transport instance.
	Connect(Transport) error

	// SendMessage sends a message to a specific network address.
	SendMessage(net.Addr, []byte) error

	// Broadcast sends a message to all connected peers.
	Broadcast([]byte) error

	// Addr returns the network address of this transport instance.
	Addr() net.Addr
}
