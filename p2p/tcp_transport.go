package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type TCPPeer struct {
	// conn is underlying connection of the peer
	conn net.Conn
	//if we dial and retrieve a conn=> outbound = true
	//if we accept and retrieve a conn=> outbound = false
	outbound bool
}

// NewTCPPeer create a new TCPPeer instance.
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

// RemoteAddr implement the Peer interface,
// which will return the remote address of the peer.
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

// Close implement the Peer interface,
// which will close the underlying connection of the peer.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// TCPTransportConfig is a configuration struct
type TCPTransportConfig struct {
	// ListenAddr is the address that the transport will listen on.
	ListenAddr string
	// HandshakeFunc is a function that will be called to perform the peer handshake.
	HandshakeFunc HandshakeFunc
	// Decoder is the decoder that will be used to decode incoming messages.
	Decoder Decoder
	// OnPeer is a function that will be called when a new peer is connected.
	OnPeer func(Peer) error
}

// TCPTransport is a transport implementation that uses TCP as the underlying transport.
type TCPTransport struct {
	TCPTransportConfig
	// listener is the underlying listener of the transport.
	listener net.Listener
	// rpcch is the channel that will be used to send incoming messages to the consumer.
	rpcch chan RPC
}

// NewTCPTransport create a new TCPTransport instance.
func NewTCPTransport(cfg TCPTransportConfig) *TCPTransport {
	return &TCPTransport{
		TCPTransportConfig: cfg,
		rpcch:              make(chan RPC),
	}
}

// Consume implement the transport interface,
// which will return read-only channel for reading incoming messages
// from another peer.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// Close implement the transport interface,
// which will close the underlying listener of the transport.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implement the transport interface,
// which will dial to the remote address.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)
	return nil
}

// ListenAndAccept implement the transport interface,
// which will start listening for incoming connections and handle them.
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	// Listen for incoming connections.
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	log.Printf("TCP Transport Listening on %s\n", t.ListenAddr)

	// Start accepting incoming connections.
	go t.StartAcceptLoop()
	return nil
}

// StartAcceptLoop will start accepting incoming connections
func (t *TCPTransport) StartAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP Accept Error: %s\n", err)
		}

		// Handle the new connection in a new goroutine.
		go t.handleConn(conn, false)
	}
}

// handleConn will handle the incoming connection.
func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	defer func() {
		fmt.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()

	// Create a new peer instance.
	peer := NewTCPPeer(conn, outbound)
	// Perform the handshake with the peer.
	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Start reading incoming messages from the peer.
	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			// fmt.Printf("TCP Read Error: %s\n", err)
			return
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}

}
