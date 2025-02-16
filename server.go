package main

import (
	"fmt"
	"log"

	"github.com/Aergiaaa/dfs/p2p"
)

type FileServerConfig struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
}

type FileServer struct {
	FileServerConfig

	store  *Store
	quitch chan struct{}
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
	}
}

func (fs *FileServer) Stop() {

	close(fs.quitch)
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

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.loop()

	return nil
}
