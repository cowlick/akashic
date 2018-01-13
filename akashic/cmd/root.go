package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

var (
	VERSION string
	version bool
)

var rootCmd = &cobra.Command{
	Use:  "akashic",
	Long: "Command-line utility for Akashic Engine",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			cmd.Println(cmd.Use + " " + VERSION)
			return
		}
		cmd.Help()
	},
}

func Execute(version string) {

	VERSION = version

	args := os.Args
	if len(args) > 1 {
		trySearchSUbCommand(args)
	}

	rootCmd.SetOutput(os.Stdout)
	err := rootCmd.Execute()
	if err != nil {
		exitError(err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, "print the version")
}

func trySearchSUbCommand(args []string) {
	subcommand := args[1]

	for _, c := range rootCmd.Commands() {
		if c.Use == subcommand {
			return
		}
	}

	path, err := findAkashicCommandPath(rootCmd.Use, subcommand)
	if err != nil {
		return
	}

	sub := &cobra.Command{
		Use:                subcommand,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := exec.Command(path, args...)
			c.Stdout = os.Stdout
			c.Stdin = os.Stdin
			c.Stderr = os.Stderr
			return c.Run()
		},
	}
	rootCmd.AddCommand(sub)
}

func exitError(err error) {
	rootCmd.SetOutput(os.Stderr)
	rootCmd.Println(err)
	os.Exit(1)
}
