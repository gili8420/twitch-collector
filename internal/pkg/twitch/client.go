package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/awend0/twitch-collector/internal/pkg/ds"
	"github.com/awend0/twitch-collector/internal/pkg/errcodes"
	"github.com/awend0/twitch-collector/internal/pkg/log"
	helix "github.com/nicklaw5/helix/v2"
)

type Client struct {
	cli *helix.Client
	cfg *Config
}

func New(cfg *Config) (*Client, error) {
	cli, err := helix.NewClient(&helix.Options{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	resp, err := cli.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		return nil, err
	}

	cli.SetAppAccessToken(resp.Data.AccessToken)
	cli.SetRefreshToken(resp.Data.RefreshToken)

	cli.OnUserAccessTokenRefreshed(func(_, _ string) {
		log.Global().Info("twitch token refreshed")
	})

	return &Client{cli: cli, cfg: cfg}, nil
}

func (c *Client) GetStreamInfo(ctx context.Context, login string) (ds.StreamInfo, error) {
	streams, err := c.cli.GetStreams(&helix.StreamsParams{
		UserLogins: []string{login},
	})
	if err != nil {
		return ds.StreamInfo{}, err
	}

	if streams.StatusCode != http.StatusOK {
		return ds.StreamInfo{}, errcodes.New(streams.StatusCode, "unexpected Twitch API answer", streams.Error, streams.ErrorMessage)
	}

	if len(streams.Data.Streams) != 1 {
		return ds.StreamInfo{}, errcodes.New(http.StatusInternalServerError, "unexpected amount of fetched streams")
	}

	return ds.StreamInfo{
		Title:           streams.Data.Streams[0].Title,
		Category:        streams.Data.Streams[0].GameName,
		Language:        streams.Data.Streams[0].Language,
		StreamStartedAt: streams.Data.Streams[0].StartedAt,
		ViewersCount:    streams.Data.Streams[0].ViewerCount,
	}, nil
}

// https://imgur.com/a/XbctLtN
func (c *Client) GetPlaybackToken(ctx context.Context, login string) (ds.PlaybackToken, error) {
	type gqlRequest struct {
		OperationName string `json:"operationName"`
		Query         string `json:"query"`
		Variables     any    `json:"variables"`
	}

	type playbackVars struct {
		Login      string `json:"login"`
		Platform   string `json:"platform"`
		PlayerType string `json:"playerType"`
	}

	type gqlResponse struct {
		Data struct {
			StreamPlaybackAccessToken struct {
				Value     string `json:"value"`
				Signature string `json:"signature"`
			} `json:"streamPlaybackAccessToken"`
		} `json:"data"`
	}

	reqBody := gqlRequest{
		OperationName: "StreamPlaybackAccessToken",
		Query: `query StreamPlaybackAccessToken(
  $login: String!,
  $platform: String!,
  $playerType: String!
) {
  streamPlaybackAccessToken(
    channelName: $login,
    params: { platform: $platform, playerBackend: "mediaplayer", playerType: $playerType }
  ) {
    value
    signature
  }
}`,
		Variables: playbackVars{
			Login:      login,
			Platform:   "web",
			PlayerType: "site",
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return ds.PlaybackToken{}, errcodes.New(http.StatusInternalServerError, "failed to marshal Twitch GQL request", err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://gql.twitch.tv/gql", bytes.NewReader(bodyBytes))
	if err != nil {
		return ds.PlaybackToken{}, errcodes.New(http.StatusInternalServerError, "failed to create Twitch GQL request", err.Error())
	}

	req.Header.Set("Client-Id", "kimne78kx3ncx6brgo4mv6wki5h1ko")
	req.Header.Set("Device-Id", "EUjL4jbzIKCH0wucs1cxBnJHDb1USUKc")
	// req.Header.Set("Origin", "https://twitch.tv")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ds.PlaybackToken{}, errcodes.New(http.StatusInternalServerError, "failed to perform Twitch GQL request", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ds.PlaybackToken{}, errcodes.New(http.StatusInternalServerError, "unexpected Twitch GQL request response", resp.Status)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ds.PlaybackToken{}, errcodes.New(http.StatusInternalServerError, "failed to read body", err.Error())
	}

	var gqlResp gqlResponse
	if err := json.Unmarshal(respBytes, &gqlResp); err != nil {
		return ds.PlaybackToken{}, errcodes.New(http.StatusInternalServerError, "failed to decode Twitch GQL request response", err.Error())
	}

	if gqlResp.Data.StreamPlaybackAccessToken.Value == "" {
		return ds.PlaybackToken{}, errcodes.New(http.StatusInternalServerError, "empty playback token")
	}

	if gqlResp.Data.StreamPlaybackAccessToken.Signature == "" {
		return ds.PlaybackToken{}, errcodes.New(http.StatusInternalServerError, "empty playback token signature")
	}

	return ds.PlaybackToken{
		Token:     gqlResp.Data.StreamPlaybackAccessToken.Value,
		Signature: gqlResp.Data.StreamPlaybackAccessToken.Signature,
	}, nil
}
