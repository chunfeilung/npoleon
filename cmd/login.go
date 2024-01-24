package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"npoleon/internal/lastfm"
	"os"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Last.fm",
	Long: `Npoleon needs permissions to scrobble tracks on your behalf. Log in to Last.fm
to provide permission to Npoleon.`,
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("LASTFM_SESSION_KEY") != "" {
			err := restoreLastFmSession()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		} else {
			err := startNewSession()
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}
		fmt.Println("Great success! You can now scrobble tracks for NPO Radio 1, 2 and 3FM.")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func restoreLastFmSession() error {
	_, err := lastfm.CreateAuthenticatedClient(
		os.Getenv("LASTFM_API_KEY"),
		os.Getenv("LASTFM_API_SECRET"),
		os.Getenv("LASTFM_SESSION_KEY"),
	)
	return err
}

func startNewSession() error {
	scrobbler, err := lastfm.CreateClient(
		os.Getenv("LASTFM_API_KEY"),
		os.Getenv("LASTFM_API_SECRET"),
	)
	if err != nil {
		return err
	}

	authUrl, token, err := scrobbler.GetAuthTokenUrl()

	if err != nil {
		return err
	}

	fmt.Println("Please follow the instructions on the page below, then return here:")
	fmt.Println(authUrl)
	fmt.Println("")
	fmt.Println("Press the Return key once you have allowed access")
	_, _ = fmt.Scanln()

	return scrobbler.Login(token)
}
