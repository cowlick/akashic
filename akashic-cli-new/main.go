package main

import (
	"fmt"
	"github.com/urfave/cli"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func npmInstall(pkg string) error {
	cmd := exec.Command("npm", "i", "-g", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func searchPackageDir(pkg string) (string, error) {
	path, err := exec.Command("npm", "root", "-g").Output()
	if err != nil {
		return "", err
	}
	return filepath.Join(strings.TrimRight(string(path), "\n"), pkg), nil
}

func copyFile(from string, to string) error {
	original, err := os.Open(from)
	if err != nil {
		return err
	}
	defer original.Close()

	target, err := os.Create(to)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, original)
	if err != nil {
		return err
	}

	return target.Sync()
}

func copyFiles(targetDir string, outDir string) error {
	return filepath.Walk(targetDir,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(targetDir, path)
			if err != nil {
				return err
			}

			return copyFile(path, filepath.Join(outDir, rel))
		})
}

func generate(pkg string, install bool) error {
	if install {
		err := npmInstall(pkg)
		if err != nil {
			return err
		}
	}

	packageDir, err := searchPackageDir(pkg)
	if err != nil {
		return err
	}

	current, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	return copyFiles(packageDir, current)
}

func main() {
	app := cli.NewApp()
	app.Name = "akashic new"
	app.Usage = "Generate project skeleton"
	app.Version = "0.0.1"

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

	app.Run(os.Args)
}
