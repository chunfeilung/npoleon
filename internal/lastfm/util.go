package lastfm

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"npoleon/internal/nporadio"
	"os"
	"strings"
)

var userHomeDir = func() (string, error) {
	return os.UserHomeDir()
}

func Initialize() {
	dir := getApplicationDir()
	_ = godotenv.Load(dir + "/config")
}

func getApplicationDir() string {
	dirname, err := userHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return dirname + "/.npoleon"
}

func appendToFile(contents string, file string) error {
	path := fmt.Sprintf("%s/%s", getApplicationDir(), file)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(contents + "\n")

	return err
}

func hasBeenScrobbled(track nporadio.Track) (bool, error) {
	date := track.PlayedAt.Format("2006-01-02")
	path := fmt.Sprintf("%s/%s.log", getApplicationDir(), date)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if line == track.PlayIdentifier() {
			return true, nil
		}
	}

	return false, nil
}

func logScrobble(track nporadio.Track) error {
	date := track.PlayedAt.Format("2006-01-02")

	return appendToFile(track.PlayIdentifier(), fmt.Sprintf("%s.log", date))
}
