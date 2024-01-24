package lastfm

import (
	"errors"
	"github.com/google/uuid"
	"npoleon/internal/nporadio"
	"os"
	"strings"
	"testing"
	"time"
)

func TestClient_GetAuthTokenUrl(t *testing.T) {
	// > Arrange
	api := &FakeApi{}
	CreateApi = func(key string, secret string) ApiInterface {
		return api
	}
	client, _ := CreateClient("whitney", "spears")

	// > Act
	url, token, _ := client.GetAuthTokenUrl()

	// > Assert
	if !strings.Contains(url, token) {
		t.Errorf("Auth token url does not contain token")
	}
}

func TestCreateClient(t *testing.T) {
	t.Run("Key and secret have not been configured yet", func(t *testing.T) {
		// > Act
		_, err := CreateClient("", "")

		// > Assert
		if err == nil {
			t.Errorf("A Last.fm client was created without credentials")
		}
	})

	t.Run("Creation of unauthenticated Last.fm client", func(t *testing.T) {
		// > Act
		client, _ := CreateClient("key", "secret")

		// > Assert
		if client == nil {
			t.Errorf("Failed to create unauthenticated Last.fm client")
		}
	})
}

func TestCreateAuthenticatedClient(t *testing.T) {
	// > Arrange
	api := &FakeApi{}
	CreateApi = func(key string, secret string) ApiInterface {
		return api
	}

	// > Act
	_, _ = CreateAuthenticatedClient("key", "secret", "d4ftpunk")

	// > Assert
	if api.SessionKey != "d4ftpunk" {
		t.Errorf("Expected session key '%v', got '%v'", "d4ftpunk", api.SessionKey)
	}
}

func TestClient_Login(t *testing.T) {
	t.Run("Login successful", func(t *testing.T) {
		// > Arrange
		dir := createTestFile(".npoleon/config", "")
		defer os.RemoveAll(dir)

		api := &FakeApi{}
		CreateApi = func(key string, secret string) ApiInterface {
			return api
		}
		client, _ := CreateAuthenticatedClient("alien", "ant", "farm")

		// > Act
		err := client.Login("token")

		// > Assert
		if err != nil {
			t.Errorf("Login should have succeeded")
			// Was session key written?
			// Need to set home dir again
		}
		contents, _ := os.ReadFile(dir + ".npoleon/config")
		if string(contents) != "LASTFM_SESSION_KEY=farm\n" {
			t.Errorf("Last.fm key not recorded, found '%v'", string(contents))
		}
	})

	t.Run("Login failed", func(t *testing.T) {
		// > Arrange
		api := &FakeApi{LoginWithTokenResult: errors.New("!")}
		CreateApi = func(key string, secret string) ApiInterface {
			return api
		}
		client, _ := CreateAuthenticatedClient("alien", "ant", "battlestar galactica")

		// > Act
		err := client.Login("token")

		// > Assert
		if err == nil {
			t.Errorf("Login should have failed")
		}
	})
}

func TestClient_Scrobble(t *testing.T) {
	dir := createTestFile(".npoleon/config", "")
	defer os.RemoveAll(dir)

	api := &FakeApi{}
	CreateApi = func(key string, secret string) ApiInterface {
		return api
	}

	t.Run("Scrobble new track", func(t *testing.T) {
		// > Arrange
		client, _ := CreateAuthenticatedClient("egg", "shaped", "head")
		track := nporadio.Track{
			Id:       uuid.New(),
			Artist:   "Fei Yu-ching",
			Title:    "一剪梅",
			PlayedAt: time.Now(),
		}

		// > Act
		err := client.Scrobble(track)

		// > Assert
		if err != nil {
			t.Errorf("failed to scrobble track %v", err)
		}
	})
}
