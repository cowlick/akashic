package main

import (
	"encoding/json"
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

func FindCommandPath(command string) (*CommandPath, error) {

	currentPath, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}
	path := filepath.Join(currentPath, "node_modules/.bin", command)
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

func NpmInstall(pkg string, global bool) error {
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

func Bootstrap(global bool) error {

	for _, pkg := range packages {

		err := NpmInstall(pkg, global)
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

func PackageVersion(pkg string) (*CommandPackageInfo, error) {

	path, err := FindCommandPath(strings.Split(pkg, "/")[1])
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

	return &CommandPackageInfo{semver.MustParse(data.Version), path.Type}, nil
}

type DistTags struct {
	Latest string `json:"latest"`
}

func UpdatePackage() error {

	for _, pkg := range packages {

		previous, err := PackageVersion(pkg)
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

		if previous.Version.LT(semver.MustParse(tags.Latest)) {
			global := false
			if previous.Type == GLOBAL {
				global = true
			}
			err = NpmInstall(pkg, global)
			if err != nil {
				return err
			}
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

		path, err := FindCommandPath(app.Name + "-cli-" + subcommand)
		if err != nil {
			return err
		}

		app.Commands = append(app.Commands, cli.Command{
			Name: subcommand,
			Action: func(c *cli.Context) error {
				cmd := exec.Command(path.Value, strings.Join(os.Args[2:], " "))
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
			Name:  "upgrade",
			Usage: "Try to update official akashic-cli-*",
			Action: func(c *cli.Context) error {
				err := UpdatePackage()
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
				return SelfUpdate(app.Version)
			},
		},
	}

	app.Run(os.Args)
}
