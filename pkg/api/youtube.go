package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/mozillazg/go-slugify"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

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

func (s *Service) GetYoutubeTrack(track *Track) {
	query := url.QueryEscape(fmt.Sprintf("%s+-+%s", slugify.Slugify(track.Name), slugify.Slugify(track.Artist)))
	request := fmt.Sprintf("https://content.googleapis.com/youtube/v3/search?q=%s&part=id,snippet&key=%s&max-results=5", query, s.Config.YoutubeKey)
	resp, err := http.Get(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var results YoutubeResults
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Fatal(err)
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
	if distance <= 20 {
		track.YoutubeUri = fmt.Sprintf("https://www.youtube.com/watch?v=%s", results.Items[bestResult].ID.VideoID)
	}
}
