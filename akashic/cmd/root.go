package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var (
	VERSION string
	version bool
	rootCmd *cobra.Command
)

func Execute(version string) {

	VERSION = version

	rootCmd.SetOutput(os.Stdout)
	err := rootCmd.Execute()
	if err != nil {
		rootCmd.SetOutput(os.Stderr)
		rootCmd.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd = &cobra.Command{
		Use:  "akashic",
		Long: "Command-line utility for Akashic Engine",
		RunE: func(cmd *cobra.Command, args []string) error {
			if version {
				cmd.Println(rootCmd.Use + " " + VERSION)
				return nil
			}

			if len(args) <= 0 {
				rootCmd.Help()
				return nil
			}

			subcommand := args[0]
			path, err := findAkashicCommandPath(rootCmd.Use, subcommand)
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

	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "print the version")
}
