package main

import (
	"github.com/awend0/twitch-collector/internal/pkg/recorder"
	"github.com/awend0/twitch-collector/internal/pkg/twitch"
	"github.com/awend0/twitch-collector/internal/repository"
	"github.com/awend0/twitch-collector/internal/service"
	"github.com/awend0/twitch-collector/internal/transport/rest"
	"github.com/awend0/twitch-collector/internal/transport/rest/handler"
	"go.uber.org/zap"
)

const (
	servicePrefix = "COLLECTOR"
)

func main() {
	lg, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	repoConfig, err := repository.NewConfig(servicePrefix)
	if err != nil {
		lg.Fatal("failed to init repo config", zap.Error(err))
	}

	repo, err := repository.New(repoConfig)
	if err != nil {
		lg.Fatal("failed to init repo", zap.Error(err))
	}

	twitchConfig, err := twitch.NewConfig(servicePrefix)
	if err != nil {
		lg.Fatal("failed to init twitch config", zap.Error(err))
	}

	twitchClient, err := twitch.New(twitchConfig)
	if err != nil {
		lg.Fatal("failed to init twitch client", zap.Error(err))
	}

	recorder := recorder.New()

	serviceConfig, err := service.NewConfig(servicePrefix)
	if err != nil {
		lg.Fatal("failed to init service config", zap.Error(err))
	}

	service := service.New(repo, twitchClient, recorder, serviceConfig)

	restHandler := handler.New(service)

	restConfig, err := rest.NewConfig(servicePrefix)
	if err != nil {
		lg.Fatal("failed to init rest config", zap.Error(err))
	}

	transport, err := rest.New(lg, restHandler, restConfig)
	if err != nil {
		lg.Fatal("failed to init rest transport", zap.Error(err))
	}

	for {
		transport.Run()

		if r := recover(); r != nil {
			lg.Warn("panic occured", zap.Any("recover", r))
		}
	}
}
