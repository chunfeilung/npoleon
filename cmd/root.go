package cmd

import (
	"npoleon/internal/lastfm"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "npoleon",
	Short: "Npoleon scrobbles tracks from NPO Radio 1, 2, and 3FM to Last.fm",
	Long: `Npoleon is a command-line utility that helps you scrobble tracks that are
being, have been, or will be played on NPO Radio 1, 2 and 3FM to Last.fm, a
social network centred around music that – like 3FM – inexplicably still exists
despite more than a decade of declining market share.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	lastfm.Initialize()
}
