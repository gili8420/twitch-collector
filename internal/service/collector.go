package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/awend0/twitch-collector/internal/pkg/ds"
	"github.com/awend0/twitch-collector/internal/pkg/log"
	"github.com/awend0/twitch-collector/internal/pkg/util"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) StartRecording(ctx context.Context, login string) (uuid.UUID, error) {
	streamInfo, err := s.twitch.GetStreamInfo(ctx, login)
	if err != nil {
		return uuid.Nil, err
	}

	playbackToken, err := s.twitch.GetPlaybackToken(ctx, login)
	if err != nil {
		return uuid.Nil, err
	}

	recordingID := uuid.New()
	log.ToCtx(ctx, log.Ctx(ctx).With(zap.String("recording_id", recordingID.String())))

	if err := s.repo.CreateRecording(ctx, ds.Recording{
		ID:              recordingID,
		Login:           login,
		CreatedAt:       time.Now(),
		Ready:           false,
		Title:           streamInfo.Title,
		Category:        streamInfo.Category,
		Language:        streamInfo.Language,
		Viewers:         streamInfo.ViewersCount,
		StreamStartedAt: streamInfo.StreamStartedAt,
	}); err != nil {
		return uuid.Nil, err
	}

	signedURL := fmt.Sprintf("https://usher.ttvnw.net/api/channel/hls/%s.m3u8?token=%s&sig=%s&allow_source=true",
		login, url.QueryEscape(playbackToken.Token), url.QueryEscape(playbackToken.Signature),
	)

	outputFilename := fmt.Sprintf("%s.mp4", recordingID)

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.RecordingDuration+30)*time.Second)
		defer cancel()

		if err := s.recorder.StartRecording(signedURL, s.cfg.RecordingDuration, outputFilename); err != nil {
			log.Ctx(ctx).Error("failed to record", zap.Error(err))
			return
		}

		log.Ctx(ctx).Info("recording finished")

		if err := s.ProcessReadyRecording(bgCtx, recordingID); err != nil {
			log.Ctx(ctx).Error("failed to process finished recording", zap.Error(err))
		}
	}()

	return recordingID, nil
}

func (s *Service) ProcessReadyRecording(ctx context.Context, recordingID uuid.UUID) error {
	if err := s.repo.UpdateRecording(ctx, recordingID, ds.UpdateRecordingRequest{
		Ready: util.ToPtr(true),
	}); err != nil {
		return err
	}

	// Planned: Upload file to S3

	return nil
}
