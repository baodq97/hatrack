//go:build !windows

package main

import (
	"os"
	"os/exec"
)

// runLogin runs claude on the terminal as is; the terminal echoes the pasted code itself.
func runLogin(c *exec.Cmd) error {
	c.Stdin = os.Stdin
	return c.Run()
}
