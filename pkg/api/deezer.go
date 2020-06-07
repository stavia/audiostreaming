package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/mozillazg/go-slugify"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// ErrSearchDeezerTrackFailed is used when a search request has failed.
var ErrSearchDeezerTrackFailed = errors.New("Search deezer track has failed")

// DeezerResults defines the json returned by Youtube
type DeezerResults struct {
	Data []struct {
		Title  string `json:"title"`
		Link   string `json:"link"`
		Artist struct {
			Name string `json:"name"`
		} `json:"artist"`
	} `json:"data"`
}

// SetDeezerURI tries to set the youtube URI of the given track
func (s *Service) SetDeezerURI(track *Track) error {
	query := url.QueryEscape(fmt.Sprintf("%s %s", slugify.Slugify(track.Name), slugify.Slugify(track.Artist)))
	request := fmt.Sprintf("http://api.deezer.com/2.0/search?q=%s&output=json", query)
	resp, err := http.Get(request)
	if err != nil {
		return ErrSearchDeezerTrackFailed
	}
	defer resp.Body.Close()
	if s.Config.LevenshteinDistance == 0 {
		s.Config.LevenshteinDistance = LevenshteinDistance
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	track.DeezerURI, err = s.GetBestDeezerResult(body, track)
	if err != nil {
		return err
	}
	return nil
}

// GetBestDeezerResult returns the best youtube results using Levenshtein distance
func (s *Service) GetBestDeezerResult(body []byte, track *Track) (uri string, err error) {
	var results DeezerResults
	err = json.Unmarshal(body, &results)
	if err != nil {
		return uri, err
	}
	artistAndSoundtrack := fmt.Sprintf("%s-%s", slugify.Slugify(track.Artist), slugify.Slugify(track.Name))
	var distance int
	var bestResult int
	for key, result := range results.Data {
		title := slugify.Slugify(result.Title)
		title = fmt.Sprintf("%s-%s", slugify.Slugify(result.Artist.Name), title)
		newDistance := levenshtein.DistanceForStrings([]rune(title), []rune(artistAndSoundtrack), levenshtein.DefaultOptions)
		if key == 0 || newDistance < distance {
			distance = newDistance
			bestResult = key
		}
	}
	if len(results.Data) > 0 && distance <= 20 {
		uri = results.Data[bestResult].Link
	}
	return uri, nil
}
