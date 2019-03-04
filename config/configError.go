package config

import "fmt"

type ConfigError struct {
	reason string
}

func (error *ConfigError) Error() string {
	return fmt.Sprintf("Configuration failed: %s", error.reason)
}
