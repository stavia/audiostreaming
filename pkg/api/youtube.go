package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/mozillazg/go-slugify"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// ErrSearchYoutubeTrackFailed is used when a search request has failed.
var ErrSearchYoutubeTrackFailed = errors.New("Search youtube track has failed")

// YoutubeResults defines the json returned by Youtube
type YoutubeResults struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

// SetYoutubeURI tries to set the youtube URI of the given track
func (s *Service) SetYoutubeURI(track *Track) error {
	query := url.QueryEscape(fmt.Sprintf("%s+-+%s", slugify.Slugify(track.Name), slugify.Slugify(track.Artist)))
	request := fmt.Sprintf("https://content.googleapis.com/youtube/v3/search?q=%s&part=id,snippet&key=%s&max-results=5", query, s.Config.YoutubeKey)
	resp, err := http.Get(request)
	if err != nil {
		return ErrSearchYoutubeTrackFailed
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	track.YoutubeURI, err = s.GetBestYoutubeResult(body, track)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// GetBestYoutubeResult returns the best youtube results using Levenshtein distance
func (s *Service) GetBestYoutubeResult(body []byte, track *Track) (uri string, err error) {
	var results YoutubeResults
	err = json.Unmarshal(body, &results)
	if err != nil {
		return uri, err
	}
	artistAndSoundtrack := fmt.Sprintf("%s-%s", slugify.Slugify(track.Artist), slugify.Slugify(track.Name))
	var distance int
	var bestResult int
	for key, result := range results.Items {
		title := slugify.Slugify(result.Snippet.Title)
		title = strings.Replace(title, "-official-audio", "", 1)
		newDistance := levenshtein.DistanceForStrings([]rune(title), []rune(artistAndSoundtrack), levenshtein.DefaultOptions)
		if key == 0 || newDistance < distance {
			distance = newDistance
			bestResult = key
		}
	}
	if len(results.Items) > 0 && distance <= 20 {
		uri = fmt.Sprintf("https://www.youtube.com/watch?v=%s", results.Items[bestResult].ID.VideoID)
	}
	return uri, nil
}
