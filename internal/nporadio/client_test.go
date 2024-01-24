package nporadio

import (
	"fmt"
	"npoleon/internal/http"
	"npoleon/internal/util"
	"os"
	"testing"
	"time"
)

func TestGetBuildId(t *testing.T) {
	httpClient := http.FakeClient{
		Responses: make(map[string][]byte),
	}

	t.Run("Web page provides buildId for API", func(t *testing.T) {
		// > Arrange
		httpClient.MakeFetchReturn(
			"https://www.nporadio2.nl/",
			`
		<!DOCTYPE html>
		<html>
		<head><shoulders><knee><toe><knee><toe>
		<script>var a={"emma":"wortelboer","buildId":"youNeedToCalmDown","cheque":"spijkerman"}</script>
		</html>
		`)

		// > Act
		res, _ := GetBuildId(httpClient, NpoRadio2)

		// > Assert
		expected := "youNeedToCalmDown"
		if res != expected {
			t.Errorf("Expected %v, got %v", expected, res)
		}
	})

	t.Run("Web page does not contain buildId", func(t *testing.T) {
		// > Arrange
		httpClient.MakeFetchReturn(
			"https://www.nporadio2.nl/",
			`
		<!DOCTYPE html>
		<title>403 Forbidden</title>
		`)

		// > Act
		_, err := GetBuildId(httpClient, NpoRadio2)

		// > Assert
		if err == nil {
			t.Errorf("Absence of buildId did not cause error")
		}
	})

	t.Run("Web page contains multiple buildIds", func(t *testing.T) {
		// > Arrange
		httpClient.MakeFetchReturn(
			"https://www.nporadio2.nl/",
			`
		<!DOCTYPE html>
		<script>{"buildId":"EvenAanMijnMoederVragen"}</script>
		<script>{"buildId":"HetRegentZonnestralen"}</script>
		`)

		// > Act
		res, _ := GetBuildId(httpClient, NpoRadio2)

		// > Assert
		expected := "EvenAanMijnMoederVragen"
		if res != expected {
			t.Errorf("Expected %v, got %v", expected, res)
		}
	})
}

// ----------------------------------------------------------------------------

func TestCreateClient(t *testing.T) {
	// > Arrange
	httpClient := http.FakeClient{Responses: make(map[string][]byte)}
	expected := "shakeItOff"
	httpClient.MakeFetchReturn(
		"https://www.nporadio1.nl/",
		`{"buildId":"`+expected+`"}`,
	)

	// > Act
	client, _ := CreateClient(httpClient, NpoRadio1)

	// > Assert
	if client.buildId != expected {
		t.Errorf("Expected %v, got %v", expected, client.buildId)
	}
}

// ----------------------------------------------------------------------------

func TestClient_FetchCurrent(t *testing.T) {
	createFakeResponseClient := func(page int) http.ClientInterface {
		httpClient := http.FakeClient{Responses: make(map[string][]byte)}
		httpClient.MakeFetchReturn("https://www.npo3fm.nl/", `{"buildId":"s3r10usR3qu3st"}`)

		fixture, _ := os.ReadFile("testdata/24-12-2023.json")
		httpClient.MakeFetchReturn(
			fmt.Sprintf("https://www.npo3fm.nl/_next/data/s3r10usR3qu3st/gedraaid/24-12-2023.json?page=%d&date=24-12-2023", page),
			string(fixture),
		)
		return httpClient
	}

	t.Run("Client fetches the right endpoint", func(t *testing.T) {
		// > Arrange
		httpClient := createFakeResponseClient(5)
		date, _ := time.Parse("2006-01-02", "2023-12-24")
		client, _ := CreateClient(httpClient, NpoRadio3)

		// > Act
		res, _ := client.fetchPage(date, 5)

		// > Assert
		if len(res) != 12 {
			t.Errorf("Expected %v, got %v", 12, len(res))
		}
	})

	t.Run("Client fetches the current track", func(t *testing.T) {
		// > Arrange
		httpClient := createFakeResponseClient(1)
		client, _ := CreateClient(httpClient, NpoRadio3)
		date, _ := util.ParseTime("2023-12-24 19:55")
		now = func() time.Time { return date.Time }

		// > Act
		res, err := client.FetchCurrent()

		// > Assert
		if res == nil || res.Title != "FELIZ NAVIDAD" {
			t.Errorf("Expected %v, got %v", "FELIZ NAVIDAD", err)
		}
	})

	t.Run("Client fetches a track that has finished playing", func(t *testing.T) {
		// > Arrange
		httpClient := createFakeResponseClient(1)
		client, _ := CreateClient(httpClient, NpoRadio3)
		date, _ := util.ParseTime("2023-12-24 22:08")
		now = func() time.Time { return date.Time }

		// > Act
		track, _ := client.FetchCurrent()

		// > Assert
		if track != nil {
			t.Errorf("No track should be playing right now")
		}
	})
}

// ----------------------------------------------------------------------------

func TestClient_FetchRange(t *testing.T) {
	createMultipleFakeResponsesClient := func(page int) http.ClientInterface {
		httpClient := http.FakeClient{Responses: make(map[string][]byte)}
		httpClient.MakeFetchReturn("https://www.npo3fm.nl/", `{"buildId":"fr4nkS1n4tr4"}`)

		fixture1, _ := os.ReadFile("testdata/5-1-2024-1.json")
		httpClient.MakeFetchReturn(
			"https://www.npo3fm.nl/_next/data/fr4nkS1n4tr4/gedraaid/5-1-2024.json?page=1&date=5-1-2024",
			string(fixture1),
		)

		fixture2, _ := os.ReadFile("testdata/6-1-2024-1.json")
		httpClient.MakeFetchReturn(
			"https://www.npo3fm.nl/_next/data/fr4nkS1n4tr4/gedraaid/6-1-2024.json?page=1&date=6-1-2024",
			string(fixture2),
		)

		fixture3, _ := os.ReadFile("testdata/6-1-2024-2.json")
		httpClient.MakeFetchReturn(
			"https://www.npo3fm.nl/_next/data/fr4nkS1n4tr4/gedraaid/6-1-2024.json?page=2&date=6-1-2024",
			string(fixture3),
		)

		fixture4, _ := os.ReadFile("testdata/6-1-2024-3.json")
		httpClient.MakeFetchReturn(
			"https://www.npo3fm.nl/_next/data/fr4nkS1n4tr4/gedraaid/6-1-2024.json?page=3&date=6-1-2024",
			string(fixture4),
		)

		fixture5, _ := os.ReadFile("testdata/31-12-2099-1.json")
		httpClient.MakeFetchReturn(
			"https://www.npo3fm.nl/_next/data/fr4nkS1n4tr4/gedraaid/31-12-2099.json?page=1&date=31-12-2099",
			string(fixture5),
		)

		return httpClient
	}

	t.Run("Client fetches tracks within single page", func(t *testing.T) {
		// > Arrange
		httpClient := createMultipleFakeResponsesClient(1)
		client, _ := CreateClient(httpClient, NpoRadio3)
		from, _ := util.ParseTime("2024-01-06 00:20")
		until, _ := util.ParseTime("2024-01-06 00:30")

		// > Act
		res, _ := client.FetchRange(from.Time, until.Time)

		// > Assert
		if len(res) != 2 {
			t.Errorf("Expected 2 tracks, got %v", len(res))
		}
	})

	t.Run("Client fetches tracks across days", func(t *testing.T) {
		// > Arrange
		httpClient := createMultipleFakeResponsesClient(1)
		client, _ := CreateClient(httpClient, NpoRadio3)
		from, _ := util.ParseTime("2024-01-05 23:55")
		until, _ := util.ParseTime("2024-01-06 00:05")

		// > Act
		res, _ := client.FetchRange(from.Time, until.Time)

		// > Assert
		if len(res) != 3 {
			t.Errorf("Expected 3 tracks, got %v", len(res))
		}
	})
}
