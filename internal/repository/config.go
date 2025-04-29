package repository

import (
	"fmt"

	env "github.com/caarlos0/env/v11"
)

type Config struct {
	ConnStr string `env:"CONNSTR" envDefault:"postgres://user:password@localhost:5432/collector?sslmode=disable"`
}

func NewConfig(servicePrefix string) (*Config, error) {
	cfg := &Config{}

	if err := env.ParseWithOptions(cfg, env.Options{
		Prefix: fmt.Sprintf("%s_REPOSITORY_", servicePrefix),
	}); err != nil {
		return nil, err
	}

	return cfg, nil
}
