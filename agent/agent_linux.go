//go:build linux
// +build linux

package main

import (
	"github.com/artemis19/viz/agent/cmd"
	godpi "github.com/mushorg/go-dpi"
)

func main() {
	//Start godpi
	godpi.Initialize()
	defer godpi.Destroy()

	// Command-line arguments
	cmd.Execute()

}
