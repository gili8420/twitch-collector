package handler

import (
	"context"
	"net/http"

	"github.com/awend0/twitch-collector/internal/pkg/errcodes"
	api "github.com/awend0/twitch-collector/pkg/rest"
)

var (
	_ api.Handler = (*Handler)(nil)
)

type Handler struct {
	svc Service
}

func New(svc Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) StartRecording(ctx context.Context, params api.StartRecordingParams) (*api.StartRecordingOK, error) {
	id, err := h.svc.StartRecording(ctx, params.Login)
	if err != nil {
		return nil, err
	}

	return &api.StartRecordingOK{ID: id}, nil
}

func (h *Handler) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	if errcode, ok := err.(*errcodes.ErrorCode); ok {
		resp := api.Error{
			Message: errcode.Message,
		}

		if errcode.Details != "" {
			resp.Details = api.NewOptString(errcode.Details)
		}

		return &api.ErrorStatusCode{
			StatusCode: errcode.StatusCode,
			Response:   resp,
		}
	}

	return &api.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.Error{
			Message: "undefined internal error",
			Details: api.NewOptString(err.Error()),
		},
	}
}
