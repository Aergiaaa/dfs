package p2p

import "net"

// Peer is a interface
// that represents the remote node.
type Peer interface {
	// Send will send the message to the peer.
	Send([]byte) error
	// RemoteAddr return the remote address of the peer.
	RemoteAddr() net.Addr
	// Close will close the connection to the peer.
	Close() error
}

// Transport represents a network transport
// that can be used to communicate between nodes.
// this can be the form of TCP, UDP, WebSocket, etc.
type Transport interface {
	// ListenAndAccept will start listening on the transport
	ListenAndAccept() error
	// Consume return a read-only channel for reading incoming messages.
	Consume() <-chan RPC
	// Close will close the transport
	Close() error
	// Dial will dial to the remote address
	Dial(string) error
}
