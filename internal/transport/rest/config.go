package rest

import (
	"fmt"

	env "github.com/caarlos0/env/v11"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8080"`
}

func NewConfig(servicePrefix string) (*Config, error) {
	cfg := &Config{}

	if err := env.ParseWithOptions(cfg, env.Options{
		Prefix: fmt.Sprintf("%s_REST_", servicePrefix),
	}); err != nil {
		return nil, err
	}

	return cfg, nil
}
