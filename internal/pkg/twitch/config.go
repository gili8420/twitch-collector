package twitch

import (
	"fmt"

	env "github.com/caarlos0/env/v11"
)

type Config struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
}

func NewConfig(servicePrefix string) (*Config, error) {
	cfg := &Config{}

	if err := env.ParseWithOptions(cfg, env.Options{
		Prefix: fmt.Sprintf("%s_TWITCH_", servicePrefix),
	}); err != nil {
		return nil, err
	}

	return cfg, nil
}
