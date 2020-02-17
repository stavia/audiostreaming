package api

// ApiProvider provides api operations.
type ApiProvider interface {
	GetYoutubeTrack(track *Track) error
	GetSpotifyTrack(track *Track) error
}

type Service struct {
	Config Config
}
