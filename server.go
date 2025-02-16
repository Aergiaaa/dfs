package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/Aergiaaa/dfs/p2p"
)

type FileServerConfig struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerConfig

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	store    *Store
	quitch   chan struct{}
}

func NewFileServer(cfg FileServerConfig) *FileServer {
	storeCfg := StoreConfig{
		Root:              cfg.StorageRoot,
		PathTransformFunc: cfg.PathTransformFunc,
	}
	return &FileServer{
		FileServerConfig: cfg,
		store:            NewStore(storeCfg),
		quitch:           make(chan struct{}),
		peers:            make(map[string]p2p.Peer),
	}
}

func (fs *FileServer) Stop() {

	close(fs.quitch)
}

func (fs *FileServer) OnPeer(p p2p.Peer) error {
	fs.peerLock.Lock()
	defer fs.peerLock.Unlock()

	fs.peers[p.RemoteAddr().String()] = p

	log.Printf("New peer connected: %s\n", p.RemoteAddr())

	return nil
}

func (fs *FileServer) loop() {
	defer func() {
		log.Println("File Server Stopped using quit action")
		fs.Transport.Close()
	}()

	for {
		select {
		case msg := <-fs.Transport.Consume():
			fmt.Println(msg)
		case <-fs.quitch:
			return
		}
	}
}

func (fs *FileServer) bootstrapNetwork() error {
	for _, addr := range fs.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			if err := fs.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}
	return nil
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if err := fs.bootstrapNetwork(); err != nil {
		return err
	}
	fs.loop()

	return nil
}
