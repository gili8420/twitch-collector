package handler

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	StartRecording(ctx context.Context, login string) (uuid.UUID, error)
}
