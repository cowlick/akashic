package main

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

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
