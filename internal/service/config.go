package service

import (
	"fmt"

	env "github.com/caarlos0/env/v11"
)

type Config struct {
	RecordingDuration int `env:"RECORDING_DURATION_SECONDS" envDefault:"60"`
}

func NewConfig(servicePrefix string) (*Config, error) {
	cfg := &Config{}

	if err := env.ParseWithOptions(cfg, env.Options{
		Prefix: fmt.Sprintf("%s_COLLECTOR_", servicePrefix),
	}); err != nil {
		return nil, err
	}

	return cfg, nil
}
