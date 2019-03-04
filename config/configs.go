package config

import (
	"github.com/bdbene/registry/server"
	"github.com/bdbene/registry/storage"
)

// Config maps configurations of all modules to be
// able to parse TOML config files
type Config struct {
	ServerConfigurations  server.ServerConfig
	StorageConfigurations storage.StorageConfig
}
