package service

import (
	"context"

	"github.com/awend0/twitch-collector/internal/pkg/ds"
	"github.com/google/uuid"
)

type TwitchClient interface {
	GetStreamInfo(ctx context.Context, login string) (ds.StreamInfo, error)
	GetPlaybackToken(ctx context.Context, login string) (ds.PlaybackToken, error)
}

type Repository interface {
	CreateRecording(ctx context.Context, recording ds.Recording) error
	UpdateRecording(ctx context.Context, recordingID uuid.UUID, req ds.UpdateRecordingRequest) error
}

type Recorder interface {
	StartRecording(playlistURL string, durationSeconds int, outputPath string) error
}
