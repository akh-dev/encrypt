package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

type Config struct {
	Service ServiceConf
	Storage StorageServiceConf
}

// DBConf - DB config

type ServiceConf struct {
	CtxTimeout int    `env:"CONTEXT_TIMEOUT" envDefault:"10"`
	Port       string `env:"LISTEN_PORT" envDefault:"8080"`
	Debug      bool   `env:"DEBUG" envDefault:"false"`
}

type StorageServiceConf struct {
	Host        string `env:"STORAGE_HOST" envDefault:"localhost"`
	Port        string `env:"STORAGE_PORT" envDefault:"8081"`
	StoreUri    string `env:"STORAGE_STORE_URI" envDefault:"/store"`
	RetrieveUri string `env:"STORAGE_RETRIEVE_URI" envDefault:"/retrieve"`
}

func Get() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(&cfg.Service); err != nil {
		return nil, errors.Wrap(err, "Failed to load Service config")
	}

	if err := env.Parse(&cfg.Storage); err != nil {
		return nil, errors.Wrap(err, "Failed to load Storage config")
	}

	return cfg, nil
}
