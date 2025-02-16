package p2p

// Peer is a interface
// that represents the remote node.
type Peer interface {
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
}
