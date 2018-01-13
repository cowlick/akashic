package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

type CommandType int

const (
	LOCAL CommandType = iota
	GLOBAL
)

type CommandPath struct {
	Type  CommandType
	Value string
}

func findCommandPath(command string) (*CommandPath, error) {

	localPath, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}
	path := filepath.Join(localPath, "node_modules/.bin", command)
	_, err = os.Stat(path)
	if err == nil {
		return &CommandPath{LOCAL, path}, nil
	}

	currentPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(filepath.Dir(currentPath), command)
	_, err = os.Stat(path)
	if err == nil {
		return &CommandPath{LOCAL, path}, nil
	}

	globalPath, err := exec.LookPath(command)
	if err != nil {
		return nil, err
	}
	return &CommandPath{GLOBAL, globalPath}, nil
}

func findAkashicCommandPath(baseName string, subcommand string) (string, error) {
	path, err := findCommandPath(baseName + "-" + subcommand)
	if err == nil {
		return path.Value, nil
	}

	path, err = findCommandPath(baseName + "-cli-" + subcommand)
	if err != nil {
		return "", errors.New("akashic command not found: " + subcommand)
	}
	return path.Value, nil
}

var packages = []string{
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
