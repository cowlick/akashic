package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func npmInstall(pkg string, global bool) error {
	var cmd *exec.Cmd
	if global {
		cmd = exec.Command("npm", "i", "-g", pkg)
	} else {
		cmd = exec.Command("npm", "i", "-D", pkg)
	}
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func bootstrap(global bool) error {

	for _, pkg := range packages {

		err := npmInstall(pkg, global)
		if err != nil {
			return err
		}
	}

	return nil
}

type Package struct {
	Version string `json:"version"`
}

type CommandPackageInfo struct {
	Version semver.Version
	Type    CommandType
}

func packageVersion(pkg string) (*CommandPackageInfo, error) {

	path, err := findCommandPath(strings.Split(pkg, "/")[1])
	if err != nil {
		return nil, err
	}

	var jsonPath string
	if path.Type == LOCAL {
		jsonPath = filepath.Join(filepath.Dir(path.Value), "..", pkg, "package.json")
	} else {
		jsonPath = filepath.Join(filepath.Dir(path.Value), "node_modules", pkg, "package.json")
	}

	bytes, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	var data Package
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	version, err := semver.Parse(data.Version)
	if err != nil {
		return nil, err
	}

	return &CommandPackageInfo{version, path.Type}, nil
}

type DistTags struct {
	Latest string `json:"latest"`
}

func updatePackage() error {

	for _, pkg := range packages {

		previous, err := packageVersion(pkg)
		if err != nil {
			return err
		}

		resp, err := http.Get(fmt.Sprintf("https://registry.npmjs.org/-/package/%s/dist-tags", url.PathEscape(pkg)))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		jsonData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var tags DistTags
		err = json.Unmarshal(jsonData, &tags)
		if err != nil {
			return err
		}

		latest, err := semver.Parse(tags.Latest)
		if err != nil {
			return err
		}

		if previous.Version.LT(latest) {
			global := false
			if previous.Type == GLOBAL {
				global = true
			}
			err = npmInstall(pkg, global)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func selfUpdate(version string) error {

	previous, err := semver.Parse(version)
	if err != nil {
		return err
	}
	latest, err := selfupdate.UpdateSelf(previous, "cowlick/akashic")
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

func export(baseName string, args cli.Args) error {
	if len(args) <= 0 {
		return errors.New("Usage: akashic export [format] [options]")
	}

	path, err := findAkashicCommandPath(baseName, "export-" + args.First())
	if err != nil {
		return err
	}

	cmd := exec.Command(path, strings.Join(os.Args[2:], " "))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func link(args cli.Args) error {

	path, err := findCommandPath("akashic-cli-install")
	if err != nil {
		return err
	}

	cmd := exec.Command(path.Value, strings.Join(append(args, "-l"), " "))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

		path, err := findAkashicCommandPath(app.Name, subcommand)
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
				err := bootstrap(c.Bool("global"))
				if err != nil {
					fmt.Println(err)
				}
				return err
			},
		},
		{
			Name:  "upgrade",
			Usage: "Try to update official akashic-cli-*",
			Action: func(c *cli.Context) error {
				err := updatePackage()
				if err != nil {
					fmt.Println(err)
				}
				return err
			},
		},
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
		{
			Name:        "export",
			Description: "Export an Akashic game",
			Usage:       "akashic export [format] [options]",
			Action: func(c *cli.Context) error {
				err := export(app.Name, c.Args())
				if err != nil {
					fmt.Println(err)
				}
				return err
			},
		},
		{
			Name: "link",
			Action: func(c *cli.Context) error {
				err := link(c.Args())
				if err != nil {
					fmt.Println(err)
				}
				return err
			},
		},
	}

	app.Run(os.Args)
}
