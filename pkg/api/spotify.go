package api

import (
	"context"
	"errors"
	"fmt"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/mozillazg/go-slugify"
	"github.com/zmb3/spotify"
)

// ErrSearchSpotifyTrackFailed is used when a search request has failed.
var ErrSearchSpotifyTrackFailed = errors.New("Search spotify track has failed")

// GetSpotifyToken returns the Spotify oauth2 token
func (s *Service) GetSpotifyToken() (token *oauth2.Token, err error) {
	config := &clientcredentials.Config{
		ClientID:     s.Config.SpotifyClientID,
		ClientSecret: s.Config.SpotifyClientSecret,
		TokenURL:     spotify.TokenURL,
	}
	token, err = config.Token(context.Background())
	if err != nil {
		log.Println("It couldn't get spotify token", err)
		return nil, err
	}
	return token, nil
}

// SetSpotifyURI tries to set the spotify URI of the given track
func (s *Service) SetSpotifyURI(track *Track, token *oauth2.Token) error {
	client := spotify.Authenticator{}.NewClient(token)
	query := fmt.Sprintf("%s %s", track.Name, track.Artist)
	results, err := client.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		return err
	}
	track.SpotifyURI, err = GetBestSpotifyResult(results, track)
	if err != nil {
		return err
	}
	return nil
}

// GetBestSpotifyResult returns the best spotify result
func GetBestSpotifyResult(results *spotify.SearchResult, track *Track) (uri string, err error) {
	if err != nil {
		return uri, ErrSearchSpotifyTrackFailed
	}
	if results.Tracks.Total == 1 {
		uri = string(results.Tracks.Tracks[0].URI)
	} else {
		for _, trackFound := range results.Tracks.Tracks {
			for _, artistFound := range trackFound.Artists {
				if slugify.Slugify(artistFound.Name) == slugify.Slugify(track.Artist) {
					return string(trackFound.URI), nil
				}
			}
		}
	}
	return uri, nil
}
