package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// GetConfigs reads configurations from the config.tml file
func GetConfigs(conf *Config) error {
	rawConfigs, err := readConfigurations()

	if err != nil {
		return err
	}

	_, err = toml.Decode(rawConfigs, conf)
	if err != nil {
		return err
	}

	return nil
}

func readConfigurations() (string, error) {
	location := os.Getenv("REGISTRY_CONFIG")
	if location == "" {
		location = "./config.tml"
	}

	log.Printf("Reading configurations from %s\n", location)
	rawConfigurations, err := ioutil.ReadFile(location)

	if err != nil {
		return "", &ConfigError{"Error reading configurations: " + err.Error()}
	}

	return string(rawConfigurations[:]), nil
}
