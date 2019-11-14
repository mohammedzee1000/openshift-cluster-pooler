package utils

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

//NewCommandMissingError creates a new error which is causes due to command missing
func NewCommandMissingError() error {
	return errors.New("no command provided")
}

//IsCommandMissingError checks if error is a CommandMissingError
func IsCommandMissingError(err error) bool {
	if strings.Contains(err.Error(), "no command provided") {
		return true
	}
	return false
}

//RunCommandWithUUID runs a command with uuid
func RunCommandWithUUID(uuid string, command string) (string, error) {
	if len(command) > 0 {
		c := exec.Command("sh", "-c", command)
		c.Env = os.Environ()
		c.Env = append(c.Env, "CLUSTER_UUID="+uuid)
		o, err := c.CombinedOutput()
		if err != nil {
			return "", err
		}
		return strings.TrimRight(string(o), "\n"), err
	}
	return "", NewCommandMissingError()
}
