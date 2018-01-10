package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
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

			rel, err := filepath.Rel(targetDir, path)
			if err != nil {
				return err
			}
			outPath := filepath.Join(outDir, rel)

			if info.IsDir() {
				if outPath == outDir {
					return nil
				}
				return os.Mkdir(outPath, os.ModeDir)
			}

			return copyFile(path, outPath)
		})
}

type Template struct {
	Path     string `json:"path"`
	GameJson string `json:"gameJson"`
}

func readTemplateInfo(packageDir string) (*Template, error) {
	bytes, err := ioutil.ReadFile(filepath.Join(packageDir, "template.json"))
	if err != nil {
		return nil, err
	}
	var template Template
	err = json.Unmarshal(bytes, &template)
	if err != nil {
		return nil, err
	}

	return &template, nil
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

	template, err := readTemplateInfo(packageDir)
	if err != nil {
		return err
	}

	current, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	return copyFiles(filepath.Join(packageDir, template.Path), current)
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
