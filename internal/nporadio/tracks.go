package nporadio

import (
	"fmt"
	"github.com/google/uuid"
	"npoleon/internal/util"
	"time"
)

type Response struct {
	PageProps PageProps `json:"pageProps"`
}

type PageProps struct {
	TrackPlays    []Play        `json:"trackPlays"`
	InitialValues InitialValues `json:"initialValues"`
	Pagination    Pagination    `json:"pagination"`
}

type Play struct {
	Id     string `json:"id"`
	Artist string `json:"artist"`
	Track  string `json:"track"`
	Time   string `json:"time"`
}

type InitialValues struct {
	Date string `json:"date"`
}

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	MaxPage     int `json:"maxPage"`
}

// ----------------------------------------------------------------------------

type Track struct {
	Id       uuid.UUID
	Artist   string
	Title    string
	PlayedAt time.Time
}

func (t Track) String() string {
	return fmt.Sprintf("%s – %s (%s)", t.Artist, t.Title, t.PlayedAt)
}

func (t Track) PlayIdentifier() string {
	return t.PlayedAt.Format("15:04") + " " + t.Id.String()
}

func (t Track) Equal(other Track) bool {
	return t.Id == other.Id &&
		t.Artist == other.Artist &&
		t.Title == other.Title &&
		t.PlayedAt.Equal(other.PlayedAt)
}

func (t Track) IsPlayedAt(moment time.Time) bool {
	start := t.PlayedAt.Add(-time.Minute)
	end := t.PlayedAt.Add(3 * time.Minute)

	return moment.After(start) && moment.Before(end)
}

// ----------------------------------------------------------------------------

type ByPlayedAt []Track

func (t ByPlayedAt) Len() int { return len(t) }

func (t ByPlayedAt) Less(i, j int) bool { return t[i].PlayedAt.Before(t[j].PlayedAt) }

func (t ByPlayedAt) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// ----------------------------------------------------------------------------

func convertResponse(response Response) ([]Track, error) {
	var tracks []Track

	for _, play := range response.PageProps.TrackPlays {
		track, err := convertTrack(play, response.PageProps.InitialValues.Date)
		if err != nil {
			return []Track{}, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func convertTrack(play Play, date string) (Track, error) {
	trackId, err := uuid.Parse(play.Id)
	if err != nil {
		return Track{}, err
	}

	playedAt, err := util.ParseTime(date + " " + play.Time)
	if err != nil {
		return Track{}, err
	}

	return Track{
		Id:       trackId,
		Artist:   play.Artist,
		Title:    play.Track,
		PlayedAt: playedAt.Time,
	}, nil
}
