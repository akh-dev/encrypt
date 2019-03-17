package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

type Config struct {
	Service ServiceConf
}

// DBConf - DB config

type ServiceConf struct {
	Port  string `env:"LISTEN_PORT" envDefault:"8081"`
	Debug bool   `env:"DEBUG" envDefault:"false"`
}

func Get() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(&cfg.Service); err != nil {
		return nil, errors.Wrap(err, "Failed to load Service config")
	}

	return cfg, nil
}
