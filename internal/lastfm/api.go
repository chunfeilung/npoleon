package lastfm

import (
	"errors"
	"github.com/shkh/lastfm-go/lastfm"
	"time"
)

// ----------------------------------------------------------------------------

type ApiInterface interface {
	GetToken() (string, error)
	GetAuthTokenUrl(token string) string
	LoginWithToken(token string) error
	GetSessionKey() string
	SetSession(sessionkey string)
	GetCorrection(artist string, title string) (lastfm.TrackGetCorrection, error)
	ScrobbleTrack(artist string, title string, timestamp time.Time) (lastfm.TrackScrobble, error)
}

// ----------------------------------------------------------------------------

type Api struct {
	api *lastfm.Api
}

// ----------------------------------------------------------------------------

func (a *Api) GetToken() (string, error) {
	return a.api.GetToken()
}

func (a *Api) GetAuthTokenUrl(token string) string {
	return a.api.GetAuthTokenUrl(token)
}

func (a *Api) LoginWithToken(token string) error {
	return a.api.LoginWithToken(token)
}

func (a *Api) GetSessionKey() string {
	return a.api.GetSessionKey()
}

func (a *Api) SetSession(sessionkey string) {
	a.api.SetSession(sessionkey)
}

func (a *Api) GetCorrection(artist string, title string) (lastfm.TrackGetCorrection, error) {
	return a.api.Track.GetCorrection(lastfm.P{
		"artist": artist,
		"track":  title,
	})
}

func (a *Api) ScrobbleTrack(artist string, title string, playedAt time.Time) (lastfm.TrackScrobble, error) {
	return a.api.Track.Scrobble(lastfm.P{
		"artist":       artist,
		"track":        title,
		"timestamp":    playedAt.Unix(),
		"chosenByUser": 0,
	})
}

// ----------------------------------------------------------------------------

type FakeApi struct {
	SessionKey           string
	LoginWithTokenResult error
}

func (f *FakeApi) GetToken() (string, error) {
	return "token", nil
}

func (f *FakeApi) GetAuthTokenUrl(token string) string {
	return "http://example.com/auth?token=" + token
}

func (f *FakeApi) LoginWithToken(token string) error {
	return f.LoginWithTokenResult
}

func (f *FakeApi) GetSessionKey() string {
	return f.SessionKey
}

func (f *FakeApi) SetSession(sessionkey string) {
	f.SessionKey = sessionkey
}

func (f *FakeApi) GetCorrection(artist string, title string) (lastfm.TrackGetCorrection, error) {
	return lastfm.TrackGetCorrection{}, errors.New("not implemented")
}

func (f *FakeApi) ScrobbleTrack(artist string, title string, playedAt time.Time) (lastfm.TrackScrobble, error) {
	return lastfm.TrackScrobble{}, nil
}

// ----------------------------------------------------------------------------

var CreateApi = func(key string, secret string) ApiInterface {
	return &Api{
		api: lastfm.New(key, secret),
	}
}
