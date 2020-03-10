package api

// Track defines the properties of a track/song to be listed
type Track struct {
	Name       string `json:"name"`
	Artist     string `json:"artist"`
	SpotifyURI string `json:"spotify_uri"`
	YoutubeURI string `json:"youtube_uri"`
	DeezerURI  string `json:"deezer_uri"`
}
