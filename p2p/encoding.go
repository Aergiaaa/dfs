package p2p

import (
	"encoding/gob"
	"io"
)

// Decoder is an interface that defines the method to decode the incoming message.
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, rpc *RPC) error {
	return gob.NewDecoder(r).Decode(rpc)
}

// DefaultDecoder is a default implementation of the Decoder interface.
type DefaultDecoder struct{}

// Decode implement the Decoder interface,
// which will read the incoming message from the reader and store it in the RPC struct.
func (dec DefaultDecoder) Decode(r io.Reader, rpc *RPC) error {
	// read the first 4 bytes to get the length of the payload
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}

	// read the rest of the payload from the reader
	rpc.Payload = buf[:n]

	return nil
}
