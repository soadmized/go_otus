package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
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
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envMap := make(Environment)

	for _, f := range files {
		path := fmt.Sprintf("%s/%s", dir, f.Name())
		value := envValue(path)
		envMap[f.Name()] = value
	}

	return envMap, nil
}

func envValue(path string) EnvValue {
	file, _ := os.OpenFile(path, os.O_RDONLY, 0o444)
	defer file.Close()

	env := EnvValue{}
	fi, _ := file.Stat()

	if fi.Size() == 0 {
		env.NeedRemove = true

		return env
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	value := scanner.Text()
	value = strings.TrimRight(value, " ")
	trimmed := bytes.ReplaceAll([]byte(value), []byte{uint8(0)}, []byte("\n"))
	env.Value = string(trimmed)

	return env
}
