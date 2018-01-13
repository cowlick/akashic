package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cowlick/akashic/npm"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func searchPackageDir(pkg string) (string, error) {
	path, err := npm.Root(true)
	if err != nil {
		return "", err
	}
	return filepath.Join(path, pkg), nil
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

// https://github.com/akashic-games/akashic-cli-commons/blob/v0.2.5/src/GameConfiguration.ts#L46
type GameConfiguration struct {
	Width             int         `json:"width"`
	Height            int         `json:"height"`
	Fps               *int        `json:"fps,omitempty"`
	Main              *string     `json:"main,omitempty"`
	Audio             interface{} `json:"audio,omitempty"`
	Assets            interface{} `json:"assets,omitempty"`
	GlobalScripts     interface{} `json:"globalScripts,omitempty"`
	OperationPlugins  interface{} `json:"operationPlugins,omitempty"`
	Environment       interface{} `json:"environment,omitempty"`
	ModuleMainScripts interface{} `json:"moduleMainScripts,omitempty"`
}

var (
	isPositiveNumber = func(input string) error {
		n, err := strconv.Atoi(input)
		if err != nil {
			return err
		} else if n < 0 {
			return errors.New("The number can not be negative!")
		}
		return nil
	}
)

func promptBasicParameter(label string, defaultValue int) (int, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: isPositiveNumber,
		Default:  strconv.Itoa(defaultValue),
	}
	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(result)
}

func promptBasicParameters(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var config GameConfiguration
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}

	width, err := promptBasicParameter("width", config.Width)
	if err != nil {
		return err
	}
	config.Width = width

	height, err := promptBasicParameter("height", config.Height)
	if err != nil {
		return err
	}
	config.Height = height

	var defaultFps int
	if config.Fps != nil {
		defaultFps = *config.Fps
	} else {
		defaultFps = 30
	}
	fps, err := promptBasicParameter("fps", defaultFps)
	if err != nil {
		return err
	}
	*config.Fps = fps

	bytes, err = json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0644)
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
	app.Version = "0.1.0"

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
