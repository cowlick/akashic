package cmd

import (
	"github.com/cowlick/akashic/npm"
	"github.com/spf13/cobra"
)

const GLOBALFLAG = "global"

var bootstrapCmd = &cobra.Command{
	Use:  "bootstrap",
	Long: "Try to install official akashic-cli-*",
	RunE: func(cmd *cobra.Command, args []string) error {

		global, err := cmd.Flags().GetBool(GLOBALFLAG)
		if err != nil {
			return err
		}

		for _, pkg := range packages {

			err := npm.Install(pkg, global)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)
	bootstrapCmd.Flags().BoolP(GLOBALFLAG, "g", false, "install the package globally")
}
