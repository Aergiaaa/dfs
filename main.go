package main

import (
	"log"

	"github.com/Aergiaaa/dfs/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportCfg := p2p.TCPTransportConfig{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportCfg)

	fileServerCfg := FileServerConfig{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	fs := NewFileServer(fileServerCfg)
	tcpTransport.OnPeer = fs.OnPeer

	return fs
}

func main() {
	fs1 := makeServer(":3000", "")
	fs2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(fs1.Start())
	}()

	if err := fs2.Start(); err != nil {
		log.Fatal(err)
	}
}
