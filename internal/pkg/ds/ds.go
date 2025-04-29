package ds

import (
	"time"

	"github.com/google/uuid"
)

type StreamInfo struct {
	Title           string
	Category        string
	Language        string
	ViewersCount    int
	StreamStartedAt time.Time
}

type PlaybackToken struct {
	Token     string
	Signature string
}

type Recording struct {
	ID              uuid.UUID
	Login           string
	CreatedAt       time.Time
	Ready           bool
	Title           string
	Category        string
	Language        string
	Viewers         int
	StreamStartedAt time.Time
}

type UpdateRecordingRequest struct {
	Login           *string
	CreatedAt       *time.Time
	Ready           *bool
	Title           *string
	Category        *string
	Language        *string
	Viewers         *int
	StreamStartedAt *time.Time
}
