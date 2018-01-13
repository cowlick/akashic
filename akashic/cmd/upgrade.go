package cmd

import (
	"encoding/json"
	"github.com/blang/semver"
	"github.com/cowlick/akashic/npm"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var upgradeCmd = &cobra.Command{
	Use:  "upgrade",
	Long: "Try to update official akashic-cli-*",
	RunE: func(cmd *cobra.Command, args []string) error {
		return updatePackages()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
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

func updatePackages() error {

	for _, pkg := range packages {

		previous, err := packageVersion(pkg)
		if err != nil {
			return err
		}

		tags, err := npm.GetDistTags(pkg)
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
			err = npm.Install(pkg, global)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
