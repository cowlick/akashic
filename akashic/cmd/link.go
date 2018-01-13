package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var linkCmd = &cobra.Command{
	Use: "link",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := findCommandPath("akashic-cli-install")
		if err != nil {
			return err
		}

		c := exec.Command(path.Value, strings.Join(append(args, "-l"), " "))
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
