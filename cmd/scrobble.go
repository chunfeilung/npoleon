package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"npoleon/internal/http"
	"npoleon/internal/lastfm"
	"npoleon/internal/nporadio"
	"npoleon/internal/scrobbling"
	"npoleon/internal/util"
	"os"
)

var scrobbleCmd = &cobra.Command{
	Use:   "scrobble STATION",
	Short: "Scrobble tracks for an NPO radio station",
	Long: `Scrobble tracks that have been, are being, or will be played on an NPO radio
station. Valid station names are "nporadio1", "nporadio2", and "npo3fm".

To scrobble a single track that’s currently being played, execute:

  npoleon scrobble 3fm --once

To keep scrobbling tracks indefinitely (at least until you terminate the
command), simply execute:

  npoleon scrobble 3fm

You can also scrobble all tracks that have been played since a particular
moment:

  npoleon scrobble 3fm --from "2024-01-20 14:30:00"

Or ask Npoleon to scrobble tracks until a specific time:

  npoleon scrobble 3fm --until "2024-01-20 20:55:00"

--from and --until can be combined to scrobble tracks for specific periods:

  npoleon scrobble 3fm --from "2024-01-20 14:30:00" --until "2024-01-20 20:55:00"
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return errors.New(`you must specify a station name, e.g. "nporadio" or "3fm"`)
		}

		if _, err := nporadio.GetStationId(args[0]); err != nil {
			return fmt.Errorf(`"%v" is not a valid station name`, args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		once, _ := cmd.Flags().GetBool("once")
		from, _ := cmd.Flags().GetString("from")
		until, _ := cmd.Flags().GetString("until")

		if os.Getenv("LASTFM_SESSION_KEY") == "" {
			err := errors.New("you are not authenticated, make sure you run `npoleon login` first")
			exitOnError(err)
		}
		lastfmClient, err := lastfm.CreateAuthenticatedClient(
			os.Getenv("LASTFM_API_KEY"),
			os.Getenv("LASTFM_API_SECRET"),
			os.Getenv("LASTFM_SESSION_KEY"),
		)
		exitOnError(err)

		radioClient, err := createRadioClient(args[0])
		exitOnError(err)

		scrobbler := scrobbling.CreateScrobbler(radioClient, lastfmClient)

		if once {
			err = scrobbler.ScrobbleOnce()
			exitOnError(err)
			return
		}

		if from == "" && until == "" {
			err = scrobbler.ScrobbleIndefinitely()
			exitOnError(err)
			return
		}

		fromTime, _ := util.ParseTimeFrom(from)
		untilTime, _ := util.ParseTimeUntil(until)

		if fromTime.After(untilTime) {
			exitOnError(errors.New("--from must be before --after"))
		}

		if from == "" && until != "" {
			err = scrobbler.ScrobbleUntil(untilTime)
			exitOnError(err)
			return
		}
		if from != "" && until == "" {
			err = scrobbler.ScrobbleFrom(fromTime)
			exitOnError(err)
			return
		}
		if from != "" && until != "" {
			err = scrobbler.ScrobblePeriod(fromTime, untilTime)
			exitOnError(err)
			return
		}
		exitOnError(errors.New("not sure what to do"))
	},
}

func init() {
	rootCmd.AddCommand(scrobbleCmd)

	scrobbleCmd.Flags().BoolP(
		"once",
		"o",
		false,
		"Only attempt to scrobble a track that is currently being played",
	)
	scrobbleCmd.Flags().StringP(
		"from",
		"f",
		"",
		"Scrobble starting from a moment in the past or future. Must be before --until",
	)
	scrobbleCmd.Flags().StringP(
		"until",
		"u",
		"",
		"Scrobble until a moment in the past or future. Must be after --from",
	)
}

func createRadioClient(stationName string) (nporadio.Client, error) {
	stationId, err := nporadio.GetStationId(stationName)
	if err != nil {
		return nporadio.Client{}, err
	}

	radioClient, err := nporadio.CreateClient(&http.Client{}, stationId)
	if err != nil {
		return nporadio.Client{}, err
	}
	return radioClient, nil
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
