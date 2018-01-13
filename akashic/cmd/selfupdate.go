package cmd

import (
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

var selfupdateCmd = &cobra.Command{
	Use:  "selfupdate",
	Long: "Try to update self via GitHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		previous, err := semver.Parse(VERSION)
		if err != nil {
			return err
		}
		latest, err := selfupdate.UpdateSelf(previous, "cowlick/akashic")
		if err != nil {
			return err
		}
		if latest.Version.Equals(previous) {
			cmd.Println("Current binary is the latest version", version)
		} else {
			cmd.Println("Successfully updated to version", latest.Version)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(selfupdateCmd)
}
