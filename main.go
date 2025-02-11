package main

import (
	"log"

	"github.com/Aergiaaa/dfs/p2p"
)

func main() {
	tcpConfig := p2p.TCPTransportConfig{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.GOBDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpConfig)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {
	// block forever
	}
}
