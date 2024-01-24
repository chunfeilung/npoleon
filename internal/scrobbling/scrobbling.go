package scrobbling

import (
	"fmt"
	"npoleon/internal/lastfm"
	"npoleon/internal/nporadio"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var now = func() time.Time { return time.Now() }

type Scrobbler struct {
	radioClient  nporadio.Client
	lastfmClient lastfm.ClientInterface
}

func CreateScrobbler(radio nporadio.Client, lastfm lastfm.ClientInterface) Scrobbler {
	return Scrobbler{
		radioClient:  radio,
		lastfmClient: lastfm,
	}
}

func (s Scrobbler) ScrobbleOnce() error {
	track, err := s.radioClient.FetchCurrent()
	if err != nil {
		return err
	}
	if track == nil {
		fmt.Println("Nothing is being played right now.")
		return nil
	}

	return s.lastfmClient.Scrobble(*track)
}

func (s Scrobbler) ScrobbleFrom(from time.Time) error {
	if from.After(now()) {
		s.waitUntil(from)
	}

	err := s.ScrobblePeriod(from, now())
	if err != nil {
		return err
	}

	return s.runUntilConditionIsMet(
		s.scrobbleCurrentTrack,
		func() bool {
			return false
		})
}

func (s Scrobbler) ScrobbleUntil(until time.Time) error {
	return s.runUntilConditionIsMet(
		s.scrobbleCurrentTrack,
		func() bool {
			return now().After(until)
		},
	)
}

func (s Scrobbler) ScrobblePeriod(from time.Time, until time.Time) error {
	if from.After(now()) {
		s.waitUntil(from)
	}

	if until.After(now()) {
		err := s.ScrobblePeriod(from, now())
		if err != nil {
			return err
		}

		err = s.ScrobbleOnce()
		if err != nil {
			return err
		}

		return s.ScrobbleUntil(until)
	}

	tracks, err := s.radioClient.FetchRange(from, until)
	if err != nil {
		return err
	}

	for idx, track := range tracks {
		if err = s.lastfmClient.Scrobble(track); err != nil {
			return err
		}
		if idx%20 == 0 {
			time.Sleep(time.Second)
		}
	}
	return nil
}

func (s Scrobbler) ScrobbleIndefinitely() error {
	return s.runUntilConditionIsMet(s.scrobbleCurrentTrack, func() bool {
		return false
	})
}

func (s Scrobbler) runUntilConditionIsMet(executeTask func() error, conditionMet func() bool) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	var remainingSleeps = 0

	for {
		select {
		case <-sig:
			os.Exit(0)
		default:
			if remainingSleeps > 0 {
				time.Sleep(500 * time.Millisecond)
				remainingSleeps = remainingSleeps - 1
				continue
			}

			if err := executeTask(); err != nil {
				return err
			}

			if conditionMet() {
				return nil
			}

			if remainingSleeps == 0 {
				remainingSleeps = 30
			}
		}
	}
}

func (s Scrobbler) waitUntil(from time.Time) {
	_ = s.runUntilConditionIsMet(
		func() error {
			return nil
		},
		func() bool {
			return !from.After(now())
		},
	)
}

func (s Scrobbler) scrobbleCurrentTrack() error {
	track, err := s.radioClient.FetchCurrent()
	if err != nil {
		return err
	}

	if track != nil {
		if err = s.lastfmClient.Scrobble(*track); err != nil {
			return err
		}
	}

	return nil
}
