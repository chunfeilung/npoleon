package lastfm

import (
	"errors"
	"fmt"
	"npoleon/internal/nporadio"
)

// ----------------------------------------------------------------------------

type ClientInterface interface {
	GetAuthTokenUrl() (string, string, error)
	Login(token string) error
	ResumeSession()
	Scrobble(track nporadio.Track) error
}

// ----------------------------------------------------------------------------

type Client struct {
	api        ApiInterface
	sessionKey string
}

// ----------------------------------------------------------------------------

func (c Client) GetAuthTokenUrl() (string, string, error) {
	token, err := c.api.GetToken()
	if err != nil {
		return "", "", err
	}
	return c.api.GetAuthTokenUrl(token), token, nil
}

func (c Client) Login(token string) error {
	err := c.api.LoginWithToken(token)

	if err != nil {
		return err
	}

	c.sessionKey = c.api.GetSessionKey()
	return appendToFile("LASTFM_SESSION_KEY="+c.sessionKey, "config")
}

func (c Client) Scrobble(track nporadio.Track) error {
	isScrobbled, err := hasBeenScrobbled(track)

	if err != nil {
		return errors.New(err.Error())
	}

	if isScrobbled {
		return nil
	}

	track = c.correctTrack(track)
	_, err = c.api.ScrobbleTrack(track.Artist, track.Title, track.PlayedAt)

	if err != nil {
		return errors.New("failed to scrobble " + track.String())
	}

	err = logScrobble(track)

	if err != nil {
		return errors.New("failed to record scrobble of " + track.String())
	}

	fmt.Println("Scrobbled", track.String())
	return nil
}

func (c Client) correctTrack(track nporadio.Track) nporadio.Track {
	res, _ := c.api.GetCorrection(track.Artist, track.Title)

	if res.Correction.TrackCorrected == "1" || res.Correction.ArtistCorrected == "1" {
		track.Artist = res.Correction.Track.Artist.Name
		track.Title = res.Correction.Track.Name
	}

	return track
}

// ----------------------------------------------------------------------------

func CreateAuthenticatedClient(key string, secret string, session string) (ClientInterface, error) {
	client := Client{
		api:        CreateApi(key, secret),
		sessionKey: session,
	}
	client.ResumeSession()
	return client, nil
}

func CreateClient(key string, secret string) (ClientInterface, error) {
	if key == "" || secret == "" {
		return nil, errors.New("please set LASTFM_API_KEY and LASTFM_API_SECRET before continuing")
	}

	return CreateAuthenticatedClient(key, secret, "")
}

func (c Client) ResumeSession() {
	if c.sessionKey != "" {
		c.api.SetSession(c.sessionKey)
	}
}

// ----------------------------------------------------------------------------

func CreateTestClient(api FakeApi) ClientInterface {
	return Client{
		api: &api,
	}
}
