package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	dir := t.TempDir()

	// Create test files
	files := map[string]string{
		"FOO":   "123",
		"BAR":   "value",
		"EMPTY": "",
		"TAB":   "	value",
		"TABR":  "	value	",
		"TERM":  "v\x00",
		"MLINE": "bar\nPLEASE IGNORE SECOND LINE",
		"W=ONG": "test",
	}

	for name, content := range files {
		err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644)
		require.NoError(t, err)
	}

	env, err := ReadDir(dir)
	require.NoError(t, err)

	expected := Environment{
		"FOO":   {Value: "123", NeedRemove: false},
		"BAR":   {Value: "value", NeedRemove: false},
		"EMPTY": {Value: "", NeedRemove: true},
		"TAB":   {Value: "\tvalue", NeedRemove: false},
		"TABR":  {Value: "\tvalue", NeedRemove: false},
		"TERM":  {Value: "v\n", NeedRemove: false},
		"MLINE": {Value: "bar", NeedRemove: false},
	}

	require.Equal(t, expected, env)
}
