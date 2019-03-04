package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// GetConfigs reads configurations from the config.tml file
func GetConfigs(conf *Config) error {
	rawConfigs, err := readConfigurations()

	if err != nil {
		return err
	}

	toml.Decode(rawConfigs, conf)
	return nil
}

func readConfigurations() (string, error) {
	rawConfigurations, err := ioutil.ReadFile("./config.tml")

	if err != nil {
		return "", &ConfigError{"Error reading configurations: " + err.Error()}
	}

	return string(rawConfigurations[:]), nil
}
