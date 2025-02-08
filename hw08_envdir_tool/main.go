package main

import (
	"fmt"
	"os"
)

// go run . "testdata/env" "/bin/bash" "testdata/echo.sh" arg1=1 arg2=2

func main() {
	// Place your code here.
	if len(os.Args) < 3 {
		fmt.Println("At least three arguments are required. Usage: go-envdir /path/to/env/dir command arg1 arg2 ...")
		os.Exit(1)
	}

	dir := os.Args[1]
	cmd := os.Args[2:]
	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	exitCode := RunCmd(cmd, env)
	os.Exit(exitCode)
}
