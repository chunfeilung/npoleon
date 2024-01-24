package lastfm

import (
	"github.com/google/uuid"
	"npoleon/internal/nporadio"
	"npoleon/internal/util"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createTestFile(path string, contents string) string {
	dir := os.TempDir() + uuid.New().String() + "/"

	_ = os.MkdirAll(dir+filepath.Dir(path), 0777)
	userHomeDir = func() (string, error) {
		return dir, nil
	}

	_, _ = os.Create(dir + path)
	_ = os.WriteFile(dir+path, []byte(contents), 0644)

	return dir
}

func TestInitialize(t *testing.T) {
	// > Arrange
	dir := createTestFile(".npoleon/config", "REEBOK_OR_NIKE=Rhythm of the Night")
	defer os.RemoveAll(dir)

	// > Act
	Initialize()

	// > Assert
	if os.Getenv("REEBOK_OR_NIKE") != "Rhythm of the Night" {
		t.Errorf("Environment variable has not been loaded")
	}
}

func TestLogScrobble(t *testing.T) {
	t.Run("First scrobble of the day", func(t *testing.T) {
		// > Arrange
		dir := createTestFile(".npoleon/config", "Not relevant for this test")
		defer os.RemoveAll(dir)
		playedAt, _ := util.ParseTime("2024-01-01 04:00:00")

		track := nporadio.Track{
			Id:       uuid.New(),
			Artist:   "Beck",
			Title:    "Beautiful Way",
			PlayedAt: playedAt.Time,
		}

		// > Act
		_ = logScrobble(track)

		// > Assert
		contents, err := os.ReadFile(dir + ".npoleon/2024-01-01.log")
		if err != nil {
			t.Errorf("Could not read log file: %v", err)
		}
		if string(contents) != track.PlayIdentifier()+"\n" {
			t.Errorf("Log file does not contain expected entry")
		}
	})

	t.Run("Second scrobble of the day", func(t *testing.T) {
		// > Arrange
		firstLine := "04:00 35363c71-bf46-4822-92c6-636448eacfae\n"
		dir := createTestFile(".npoleon/2024-01-01.log", firstLine)
		defer os.RemoveAll(dir)
		playedAt, _ := util.ParseTime("2024-01-01 04:04:00")

		track := nporadio.Track{
			Id:       uuid.New(),
			Artist:   "Madonna",
			Title:    "Ray of Light",
			PlayedAt: playedAt.Time,
		}

		// > Act
		_ = logScrobble(track)

		// > Assert
		contents, err := os.ReadFile(dir + ".npoleon/2024-01-01.log")
		if err != nil {
			t.Errorf("Could not read log file: %v", err)
		}
		if string(contents) != firstLine+track.PlayIdentifier()+"\n" {
			t.Errorf("Log file does not contain expected entry")
		}
	})
}

func TestHasBeenScrobbled(t *testing.T) {
	trackId := uuid.New()
	firstLine := "08:20 " + trackId.String() + "\n"
	dir := createTestFile(".npoleon/2024-01-01.log", firstLine)
	defer os.RemoveAll(dir)
	playedAt, _ := util.ParseTime("2024-01-01 08:20:00")

	t.Run("This track play has already been scrobbled", func(t *testing.T) {
		// > Arrange
		track := nporadio.Track{
			Id:       trackId,
			Artist:   "Bula Matari",
			Title:    "Taxi Drivers",
			PlayedAt: playedAt.Time,
		}

		// > Act
		res, err := hasBeenScrobbled(track)

		// > Assert
		if err != nil {
			t.Errorf("Could not determine if track has been scrobbled: %v", err)
		}
		if !res {
			t.Errorf("Track was not seen as scrobbled")
		}
	})

	t.Run("Track is played for second time on same day", func(t *testing.T) {
		// > Arrange
		track := nporadio.Track{
			Id:       trackId,
			Artist:   "Bula Matari",
			Title:    "Taxi Drivers",
			PlayedAt: playedAt.Time.Add(6 * time.Hour),
		}

		// > Act
		res, err := hasBeenScrobbled(track)

		// > Assert
		if err != nil {
			t.Errorf("Could not determine if track has been scrobbled: %v", err)
		}
		if res {
			t.Errorf("Track is seen as scrobbled")
		}
	})

	t.Run("Track has not yet been scrobbled", func(t *testing.T) {
		// > Arrange
		track := nporadio.Track{
			Id:       uuid.New(),
			Artist:   "MNDR",
			Title:    "Lock & Load (feat. Killer Mike)",
			PlayedAt: playedAt.Time.Add(7 * time.Minute),
		}

		// > Act
		res, err := hasBeenScrobbled(track)

		// > Assert
		if err != nil {
			t.Errorf("Could not determine if track has been scrobbled: %v", err)
		}
		if res {
			t.Errorf("Track is seen as scrobbled")
		}
	})
}
