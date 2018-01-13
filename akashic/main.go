package main

import "github.com/cowlick/akashic/akashic/cmd"

var (
	VERSION = "0.2.0"
)

func main() {
	cmd.Execute(VERSION)
}
