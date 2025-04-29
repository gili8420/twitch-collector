package repository

import (
	"context"

	"github.com/awend0/twitch-collector/internal/repository/sqlc"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	queries *sqlc.Queries
	cfg     *Config
}

func New(cfg *Config) (*Repository, error) {
	conn, err := pgx.Connect(context.Background(), cfg.ConnStr)
	if err != nil {
		return nil, err
	}

	return &Repository{
		queries: sqlc.New(conn),
		cfg:     cfg,
	}, nil
}
