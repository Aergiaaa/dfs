package p2p

// Peer is a interface 
// that represents the remote node. 
type Peer interface {

}

// Transport represents a network transport 
// that can be used to communicate between nodes.
// this can be the form of TCP, UDP, WebSocket, etc.
type Transport interface {
  ListenAndAccept() error
}