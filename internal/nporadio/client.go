package nporadio

import (
	"encoding/json"
	"errors"
	"fmt"
	"npoleon/internal/http"
	"regexp"
	"sort"
	"time"
)

var now = func() time.Time { return time.Now() }

func GetBuildId(httpClient http.ClientInterface, stationId StationId) (string, error) {
	mainPage := fmt.Sprintf("https://www.%s.nl/", stationId)
	resp, err := httpClient.Fetch(mainPage)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`"buildId":"([^"]+)"`)
	matches := re.FindStringSubmatch(string(resp))

	if len(matches) < 1 {
		msg := fmt.Sprintf("failed to identify buildId for %s", stationId)
		return "", errors.New(msg)
	}

	return matches[1], nil
}

type Client struct {
	httpClient http.ClientInterface
	stationId  StationId
	buildId    string
}

func CreateClient(httpClient http.ClientInterface, stationId StationId) (Client, error) {
	buildId, err := GetBuildId(httpClient, stationId)
	if err != nil {
		return Client{}, err
	}

	return Client{
		httpClient: httpClient,
		stationId:  stationId,
		buildId:    buildId,
	}, nil
}

func (c Client) fetchPage(date time.Time, page int) ([]Track, error) {
	endpoint := fmt.Sprintf(
		"https://www.%s.nl/_next/data/%s/gedraaid/%s.json?page=%d&date=%s",
		c.stationId,
		c.buildId,
		date.Format("2-1-2006"),
		page,
		date.Format("2-1-2006"),
	)

	resp, err := c.httpClient.Fetch(endpoint)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, err
	}

	tracks, err := convertResponse(response)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (c Client) FetchCurrent() (*Track, error) {
	location, _ := time.LoadLocation("Europe/Amsterdam")
	tracks, err := c.fetchPage(now().In(location), 1)

	if err != nil {
		return nil, err
	}

	if len(tracks) == 0 {
		return nil, nil
	}

	var track = tracks[0]
	if !track.IsPlayedAt(now()) {
		return nil, nil
	}

	return &track, nil
}

func (c Client) FetchRange(from time.Time, until time.Time) ([]Track, error) {
	var allTracks []Track
	var page = 1
	var now = until
	for {
		if now.Before(from) {
			break
		}

		currentTracks, err := c.fetchPage(now, page)

		if err != nil {
			return nil, err
		}

		if len(currentTracks) == 0 {
			// If 'from' time has not been reached yet, go the previous day
			break
		}

		allTracks = append(allTracks, currentTracks...)

		if containsTracksBeforeDate(currentTracks, from) {
			break
		}

		page++
		now = currentTracks[len(currentTracks)-1].PlayedAt
	}

	if from.Before(now) {
		year, month, day := now.Add(-24 * time.Hour).Date()
		location, _ := time.LoadLocation("Europe/Amsterdam")
		yesterday := time.Date(year, month, day, 23, 59, 59, 0, location)
		olderTracks, err := c.FetchRange(from, yesterday)
		if err != nil {
			return nil, err
		}
		allTracks = append(allTracks, olderTracks...)
	}

	filteredTracks := removeTracksOutsideRange(allTracks, from, until)
	sort.Sort(ByPlayedAt(filteredTracks))
	return filteredTracks, nil
}

func removeTracksOutsideRange(tracks []Track, start time.Time, end time.Time) []Track {
	var result []Track
	for _, t := range tracks {
		if t.PlayedAt.Before(start) {
			continue
		}
		if t.PlayedAt.After(end) {
			continue
		}
		result = append(result, t)
	}
	return result
}

func containsTracksBeforeDate(tracks []Track, date time.Time) bool {
	return tracks[len(tracks)-1].PlayedAt.Before(date)
}
