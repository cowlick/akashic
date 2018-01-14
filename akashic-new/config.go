package main

import (
	"encoding/json"
	"errors"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

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
