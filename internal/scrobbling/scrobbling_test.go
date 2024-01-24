package scrobbling

import (
	"fmt"
	"npoleon/internal/http"
	"npoleon/internal/lastfm"
	"npoleon/internal/nporadio"
	"testing"
	"time"
)

func createFakeHttpClient() http.ClientInterface {
	httpClient := http.FakeClient{Responses: make(map[string][]byte)}
	httpClient.MakeFetchReturn(
		"https://www.npo3fm.nl/",
		`{"buildId":"buildId"}`,
	)

	loc, _ := time.LoadLocation("Europe/Amsterdam")
	date := time.Now().In(loc).Format("2-1-2006")
	httpClient.MakeFetchReturn(
		fmt.Sprintf("https://www.npo3fm.nl/_next/data/buildId/gedraaid/%s.json?page=1&date=%s", date, date),
		"{}",
	)

	return httpClient
}

func TestScrobbler_ScrobbleOnce(t *testing.T) {
	// > Arrange
	radioClient, _ := nporadio.CreateClient(createFakeHttpClient(), nporadio.NpoRadio3)
	lastfmClient := lastfm.CreateTestClient(lastfm.FakeApi{})
	scrobbler := CreateScrobbler(radioClient, lastfmClient)

	// > Act
	err := scrobbler.ScrobbleOnce()

	// > Assert
	if err != nil {
		t.Errorf("Scrobbling failed")
	}
}
