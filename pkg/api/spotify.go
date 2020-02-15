package api

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/mozillazg/go-slugify"
	"github.com/zmb3/spotify"
)

func (s *Service) GetSpotifyTrack(track *Track) {
	config := &clientcredentials.Config{
		ClientID:     s.Config.SpotifyClientID,
		ClientSecret: s.Config.SpotifyClientSecret,
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)
	query := fmt.Sprintf("%s artist:%s", track.Name, track.Artist)
	results, err := client.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	if results.Tracks.Total == 1 {
		track.SpotifyUri = results.Tracks.Tracks[0].Endpoint
	} else {
		for _, trackFound := range results.Tracks.Tracks {
			for _, artistFound := range trackFound.Artists {
				if slugify.Slugify(artistFound.Name) == slugify.Slugify(track.Artist) {
					track.SpotifyUri = trackFound.Endpoint
					return
				}
			}
		}
	}
}
