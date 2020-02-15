package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stavia/audiostreaming/pkg/api"
)

func main() {
	arg1 := flag.String("track", "", "Name of the track you want to search")
	arg2 := flag.String("artist", "", "Artist of the track you want to search")
	flag.Parse()
	if *arg1 == "" || *arg2 == "" {
		flag.PrintDefaults()
		return
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var track api.Track
	track.Name = *arg1
	track.Artist = *arg2
	api := new(api.Service)
	api.Config.SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	api.Config.SpotifyClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	api.Config.YoutubeKey = os.Getenv("YOTUBE_KEY")

	api.GetYoutubeTrack(&track)
	api.GetSpotifyTrack(&track)
	fmt.Println(track)
}
