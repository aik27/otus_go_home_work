package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	env := Environment{}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = f.Close()
			if err != nil {
				log.Printf("Failed to close file %s: %v\n", filePath, err)
			}
		}()

		if strings.Contains(file.Name(), "=") {
			continue
		}

		fi, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}

		if fi.Size() == 0 {
			env[file.Name()] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
		} else {
			scanner := bufio.NewScanner(f)
			if scanner.Scan() {
				value := strings.TrimRight(scanner.Text(), " ")
				value = strings.TrimRight(value, "\t")
				value = strings.ReplaceAll(value, "\x00", "\n")
				env[file.Name()] = EnvValue{
					Value:      value,
					NeedRemove: false,
				}
			}

			if err := scanner.Err(); err != nil {
				return nil, err
			}
		}
	}
	return env, nil
}
