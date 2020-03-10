package api

// Config describes IDs and keys of the audio streaming applications
type Config struct {
	// SpotifyClientID is the spotify application's ID.
	SpotifyClientID string

	// SpotifyClientSecret is the spotify application's secret.
	SpotifyClientSecret string

	// YoutubeKey is the youtube key
	YoutubeKey string
}
