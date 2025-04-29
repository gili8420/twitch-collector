package rest

import (
	"fmt"
	"net/http"

	handler "github.com/awend0/twitch-collector/internal/transport/rest/handler"
	"github.com/awend0/twitch-collector/internal/transport/rest/middleware"
	api "github.com/awend0/twitch-collector/pkg/rest"
	"go.uber.org/zap"
)

type Transport struct {
	lg  *zap.Logger
	cfg *Config
	srv *api.Server
}

func New(lg *zap.Logger, handler *handler.Handler, cfg *Config) (*Transport, error) {
	srv, err := api.NewServer(handler, api.WithMiddleware(middleware.Logging(lg)))
	if err != nil {
		return nil, err
	}

	return &Transport{
		lg:  lg,
		srv: srv,
		cfg: cfg,
	}, nil
}

func (t *Transport) Run() {
	t.lg.Info("started")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", t.cfg.Port), t.srv); err != nil {
		panic(err)
	}
}
