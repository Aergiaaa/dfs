package main

import (
	"fmt"
	"log"

	"github.com/Aergiaaa/dfs/p2p"
)

func OnPeer(p p2p.Peer) error {
	// fmt.Println("doing some logic with the peer outside of the TCPtransport")
	p.Close()
	return nil
}

func main() {
	tcpConfig := p2p.TCPTransportConfig{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpConfig)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Success: %+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {
	// block forever
	}
}
