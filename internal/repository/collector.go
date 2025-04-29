package repository

import (
	"context"
	"errors"
	"net/http"

	"github.com/awend0/twitch-collector/internal/pkg/ds"
	"github.com/awend0/twitch-collector/internal/pkg/errcodes"
	"github.com/awend0/twitch-collector/internal/repository/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateRecording(ctx context.Context, recording ds.Recording) error {
	if err := r.queries.CreateRecording(ctx, sqlc.CreateRecordingParams{
		ID:              UUIDToPG(recording.ID),
		Login:           recording.Login,
		CreatedAt:       TimeToPG(recording.CreatedAt),
		Ready:           recording.Ready,
		Title:           recording.Title,
		Category:        recording.Category,
		Language:        recording.Language,
		Viewers:         int32(recording.Viewers),
		StreamStartedAt: TimeToPG(recording.StreamStartedAt),
	}); err != nil {
		return errcodes.New(http.StatusInternalServerError, "failed to create recording", err.Error())
	}

	return nil
}

func (r *Repository) UpdateRecording(ctx context.Context, recordingID uuid.UUID, req ds.UpdateRecordingRequest) error {
	if _, err := r.queries.UpdateRecording(ctx, sqlc.UpdateRecordingParams{
		RecordingID: UUIDToPG(recordingID),

		Login:           NilStringToPG(req.Login),
		CreatedAt:       NilTimeToPG(req.CreatedAt),
		Ready:           NilBoolToPG(req.Ready),
		Title:           NilStringToPG(req.Title),
		Category:        NilStringToPG(req.Category),
		Language:        NilStringToPG(req.Language),
		Viewers:         NilIntToPG(req.Viewers),
		StreamStartedAt: NilTimeToPG(req.StreamStartedAt),
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errcodes.New(http.StatusNotFound, "unknown recording")
		}

		return errcodes.New(http.StatusInternalServerError, "failed to update recording", err.Error())
	}

	return nil
}
