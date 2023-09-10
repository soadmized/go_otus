package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	err := setEnvs(env)
	if err != nil {
		log.Fatal(err)
	}

	path, err := exec.LookPath(cmd[0])
	if err != nil {
		return 2
	}

	command := exec.Command(path, cmd[1:]...)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	err = command.Run()
	if err != nil {
		exit := err.(*exec.ExitError) //nolint:errorlint

		return exit.ExitCode()
	}

	return 0
}

func setEnvs(envs Environment) error {
	for k, v := range envs {
		if v.NeedRemove {
			err := os.Unsetenv(k)
			if err != nil {
				return err
			}
		}

		err := os.Setenv(k, v.Value)
		if err != nil {
			return err
		}
	}

	return nil
}
