package main

import (
	"fmt"
	"github.com/elastic/go-ucfg/yaml"
	"inverse_time_lag_relay/curve"
	"inverse_time_lag_relay/relay"
	"os"
)

var config *Config

type Config struct {
	Relay relay.Relay `config:"relay"`
	Curve curve.Curve `config:"curve"`
}

func ParseConfiguration(filename string) error {
	configContent, err := yaml.NewConfigWithFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file not found")
		}
		return err
	}
	if err := configContent.Unpack(&config); err != nil {
		return err
	}
	return nil
}
