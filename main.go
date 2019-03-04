package main

import (
	"log"

	"github.com/bdbene/registry/config"
	"github.com/bdbene/registry/server"
	"github.com/bdbene/registry/storage"
)

// import "github.com/etcd-io/etcd/raft"

func main() {
	configs := new(config.Config)
	err := config.GetConfigs(configs)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	storage, err := storage.NewSqlStore(&configs.StorageConfigurations)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer storage.Close()

	server, err := server.NewServer(&configs.ServerConfigurations, storage)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	server.Listen()
}
