package api

// Track defines the properties of a track/song to be listed
type Track struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	SpotifyUri    string `json:"uri"`
	YoutubeUri    string `json:"uri"`
}
