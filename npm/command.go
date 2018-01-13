package npm

import (
	"os"
	"os/exec"
	"strings"
)

func Install(pkg string, global bool) error {
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

func Root(global bool) (string, error) {
	var cmd *exec.Cmd
	if global {
		cmd = exec.Command("npm", "root", "-g")
	} else {
		cmd = exec.Command("npm", "root")
	}
	path, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(path), "\n"), nil
}
