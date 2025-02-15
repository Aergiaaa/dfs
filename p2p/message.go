package p2p

import "net"

// Message hold any arbitrary data that
// is being sent over each transport
// between two nodes.
type RPC struct {
	// From is the address of the sender.
	From net.Addr
	// Payload is the data being sent.
	Payload []byte
}
