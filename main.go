package main

import (
	"log"
	"time"

	"github.com/Aergiaaa/dfs/p2p"
)

func main() {
	tcpTransportCfg := p2p.TCPTransportConfig{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// OnPeer:        OnPeer,
		// TODO: add OnPeer
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportCfg)

	fileServerCfg := FileServerConfig{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	fs := NewFileServer(fileServerCfg)

	go func() {
		time.Sleep(time.Second * 2)
		fs.Stop()
	}()

	if err := fs.Start(); err != nil {
		log.Fatal(err)
	}

}
