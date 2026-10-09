package main

import (
	"os"
	"os/exec"

	"golang.org/x/term"
)

// runLogin runs claude, masking what the user types: claude's Windows code prompt echoes
// nothing, so a pasted code would give no sign it arrived.
func runLogin(c *exec.Cmd) error {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		c.Stdin = os.Stdin
		return c.Run()
	}
	in, err := c.StdinPipe()
	if err != nil {
		return err
	}
	if err := c.Start(); err != nil {
		return err
	}
	old, err := term.MakeRaw(fd)
	if err != nil {
		c.Process.Kill()
		c.Wait()
		return err
	}
	defer term.Restore(fd, old)
	go readMasked(os.Stdin, os.Stdout, in, func() { c.Process.Kill() })
	return c.Wait()
}
