package main

import (
	"fmt"
	"github.com/cowlick/akashic/npm"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
)

func searchPackageDir(pkg string) (string, error) {
	path, err := npm.Root(true)
	if err != nil {
		return "", err
	}
	return filepath.Join(path, pkg), nil
}

func generate(pkg string, install bool) error {
	if install {
		err := npm.Install(pkg, true)
		if err != nil {
			return err
		}
	}

	packageDir, err := searchPackageDir(pkg)
	if err != nil {
		return err
	}

	template, err := readTemplateInfo(packageDir)
	if err != nil {
		return err
	}

	current, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	err = copyFiles(filepath.Join(packageDir, template.Path), current)
	if err != nil {
		return err
	}

	return promptBasicParameters(filepath.Join(current, template.GameJson))
}

func main() {
	app := cli.NewApp()
	app.Name = "akashic new"
	app.Usage = "Generate project skeleton"
	app.Version = "0.2.0"

	app.ArgsUsage = "[npm package]"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "install, i",
			Usage: "Install npm package from npm registory before generate template",
		},
	}

	app.Action = func(c *cli.Context) error {

		args := c.Args()
		if len(args) <= 0 {
			cli.ShowAppHelp(c)
			return nil
		}

		err := generate(args.First(), c.Bool("install"))
		if err != nil {
			fmt.Print(err)
		}
		return err
	}

	app.Commands = []cli.Command{
		{
			Name:  "selfupdate",
			Usage: "Try to update self via GitHub",
			Action: func(c *cli.Context) error {
				err := selfUpdate(app.Version)
				if err != nil {
					fmt.Println(err)
				}
				return err
			},
		},
	}

	app.Run(os.Args)
}
