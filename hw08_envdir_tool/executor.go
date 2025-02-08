package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	for _, arg := range cmd {
		if strings.Contains(arg, ";") || strings.Contains(arg, "&") || strings.Contains(arg, "|") {
			return 1
		}
	}

	// Place your code here.
	for k, v := range env {
		if v.NeedRemove {
			err := os.Unsetenv(k)
			if err != nil {
				return 1
			}
			continue
		}

		err := os.Setenv(k, v.Value)
		if err != nil {
			return 1
		}
	}

	absCmd, err := filepath.Abs(cmd[0])
	if err != nil {
		return 1
	}

	c := exec.Command(absCmd, cmd[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin

	err = c.Run()
	if err != nil {
		return 1
	}

	return c.ProcessState.ExitCode()
}
