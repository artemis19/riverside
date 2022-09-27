//go:build windows
// +build windows

package main

import (
	"github.com/artemis19/viz/agent/cmd"
)

func main() {
	// Command-line arguments
	cmd.Execute()
}
