package api

import (
	"context"
	"errors"
	"fmt"
	"log"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/mozillazg/go-slugify"
	"github.com/zmb3/spotify"
)

// ErrSearchSpotifyTrackFailed is used when a search request has failed.
var ErrSearchSpotifyTrackFailed = errors.New("Search spotify track has failed")

// GetSpotifyTrack tries to set the spotify URI of the given track
func (s *Service) GetSpotifyTrack(track *Track) error {
	config := &clientcredentials.Config{
		ClientID:     s.Config.SpotifyClientID,
		ClientSecret: s.Config.SpotifyClientSecret,
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("It couldn't get spotify token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)
	query := fmt.Sprintf("%s artist:%s", track.Name, track.Artist)
	results, err := client.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		return ErrSearchSpotifyTrackFailed
	}

	if results.Tracks.Total == 1 {
		track.SpotifyUri = results.Tracks.Tracks[0].Endpoint
	} else {
		for _, trackFound := range results.Tracks.Tracks {
			for _, artistFound := range trackFound.Artists {
				if slugify.Slugify(artistFound.Name) == slugify.Slugify(track.Artist) {
					track.SpotifyUri = trackFound.Endpoint
					return nil
				}
			}
		}
	}
	return nil
}
