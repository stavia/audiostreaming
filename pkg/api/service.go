package api

// ApiProvider provides api operations.
type ApiProvider interface {
	GetYoutubeTrack(track *Track)
	GetSpotifyTrack(track *Track)
}

type Service struct {
	Config Config
}
