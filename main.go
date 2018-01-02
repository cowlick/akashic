package main

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func SubCommandPath(subcommand string) (string, error) {

	// local command
	currentPath, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	path := filepath.Join(currentPath, "node_modules/.bin", subcommand)
	files, err := filepath.Glob(path)
	if err != nil {
		return "", err
	}
	if len(files) != 0 {
		return files[0], nil
	}

	// global command
	return exec.LookPath(subcommand)
}

var packages = []string{
	"@akashic/akashic-cli-commons",
	"@akashic/akashic-cli-init",
	"@akashic/akashic-cli-scan",
	"@akashic/akashic-cli-modify",
	"@akashic/akashic-cli-update",
	"@akashic/akashic-cli-install",
	"@akashic/akashic-cli-uninstall",
	"@akashic/akashic-cli-config",
	"@akashic/akashic-cli-export-html",
	"@akashic/akashic-cli-export-zip",
	"@akashic/akashic-cli-stat",
}

func Bootstrap(global bool) error {

	for _, pkg := range packages {

		path, err := exec.LookPath("npm")
		if err != nil {
			return err
		}

		args := pkg
		if global {
			args = "-g " + args
		}

		cmd := exec.Command(path + " i", args)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func SelfUpdate(version string) error {

	previous := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(previous, "cowlick/akashic-cli")
	if err != nil {
		return err
	}
	if latest.Version.Equals(previous) {
		fmt.Println("Current binary is the latest version", version)
	} else {
		fmt.Println("Successfully updated to version", latest.Version)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "akashic"
	app.Usage = "Command-line utility for Akashic Engine"
	app.Version = "0.0.1"

	app.Before = func(c *cli.Context) error {
		args := c.Args()
		if len(args) <= 0 {
			return nil
		}

		subcommand := args.First()
		for _, c := range app.Commands {
			if c.HasName(subcommand) {
				return nil
			}
		}

		path, err := SubCommandPath(app.Name + "-cli-" + subcommand)
		if err != nil {
			return err
		}

		app.Commands = append(app.Commands, cli.Command{
			Name: subcommand,
			Action: func(c *cli.Context) error {
				cmd := exec.Command(path, strings.Join(os.Args[2:], " "))
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				return cmd.Run()
			},
		})

		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "bootstrap",
			Usage: "Try to install official akashic-cli-*",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "global, g",
				},
			},
			Action: func(c *cli.Context) error {
				return Bootstrap(c.Bool("global"))
			},
		},
		{
			Name:  "selfupdate",
			Usage: "Try to update self via GitHub",
			Action: func(c *cli.Context) error {
				return SelfUpdate(app.Version)
			},
		},
	}

	app.Run(os.Args)
}
