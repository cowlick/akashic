package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var exportCmd = &cobra.Command{
	Use:  "export",
	Long: "Export an Akashic game",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) <= 0 {
			return errors.New("Usage: akashic export [format] [options]")
		}

		path, err := findAkashicCommandPath(rootCmd.Use, "export-"+args[0])
		if err != nil {
			return err
		}

		c := exec.Command(path, args[1:]...)
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		return c.Run()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
