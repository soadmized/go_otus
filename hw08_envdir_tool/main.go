package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	path := args[1]
	cmdAndArgs := args[2:]

	envs, err := ReadDir(path)
	if err != nil {
		log.Println(err)
	}

	RunCmd(cmdAndArgs, envs)

	log.Print(path, cmdAndArgs)
}
