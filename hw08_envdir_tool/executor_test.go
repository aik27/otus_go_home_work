package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
		"FOO": {Value: "123", NeedRemove: false},
		"BAR": {Value: "value", NeedRemove: false},
	}

	cmd := []string{"/bin/sh", "-c", "echo $FOO $BAR"}

	exitCode := RunCmd(cmd, env)
	require.Equal(t, 0, exitCode)

	// Verify the output
	output, err := exec.Command("/bin/sh", "-c", "echo $FOO $BAR").Output()
	require.NoError(t, err)
	require.Equal(t, "123 value\n", string(output))
}
