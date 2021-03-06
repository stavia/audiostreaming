package api

import "golang.org/x/oauth2"

// Provider provides api operations.
type Provider interface {
	SetYoutubeURI(track *Track) error
	GetBestYoutubeResult(body []byte, track *Track) (uri string, err error)
	GetSpotifyToken() (token *oauth2.Token, err error)
	SetSpotifyURI(track *Track) error
	GetBestSpotifyResult(body []byte, track *Track) (uri string, err error)
	SetDeezerURI(track *Track) error
	GetBestDeezerResult(body []byte, track *Track) (uri string, err error)
}

// Service provides setting URIs into Track and methods for searching
// into audio streaming APIs
type Service struct {
	Config Config
}
