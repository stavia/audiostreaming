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

// SetSpotifyURI tries to set the spotify URI of the given track
func (s *Service) SetSpotifyURI(track *Track) error {
	config := &clientcredentials.Config{
		ClientID:     s.Config.SpotifyClientID,
		ClientSecret: s.Config.SpotifyClientSecret,
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Println("It couldn't get spotify token", err)
		return err
	}

	client := spotify.Authenticator{}.NewClient(token)
	query := fmt.Sprintf("%s %s", track.Name, track.Artist)
	results, err := client.Search(query, spotify.SearchTypeTrack)
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
					return trackFound.Endpoint, nil
				}
			}
		}
	}
	return uri, nil
}
