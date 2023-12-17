package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type Config struct {
	LogLevel string `json:"logLevel"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	GRPCPort string `json:"grpcPort"`
}

func Load() (*Config, error) {
	var conf Config

	raw, err := os.ReadFile("config/config.json")
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = json.Unmarshal(raw, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal config")
	}

	return &conf, nil
}

func (c Config) Addr() string {
	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)

	return addr
}

func (c Config) GRPCAddr() string {
	addr := fmt.Sprintf("%s:%s", c.Host, c.GRPCPort)

	return addr
}
