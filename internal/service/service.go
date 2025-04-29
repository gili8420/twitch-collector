package service

type Service struct {
	repo     Repository
	twitch   TwitchClient
	recorder Recorder
	cfg      *Config
}

func New(repo Repository, twitch TwitchClient, recorder Recorder, cfg *Config) *Service {
	return &Service{
		repo:     repo,
		twitch:   twitch,
		recorder: recorder,
		cfg:      cfg,
	}
}
