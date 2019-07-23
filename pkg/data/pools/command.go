package pools

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func NewCommandMissingError() error {
	return errors.New("no command provided")
}

func IsCommandMissingError(err error) bool {
	if strings.Contains(err.Error(), "no command provided") {
		return true
	}
	return false
}

func runCommand(uuid string, command string) (string, error) {
	if len(command) > 0 {
		c := exec.Command("sh", "-c", command)
		c.Env = os.Environ()
		c.Env = append(c.Env, "CLUSTER_UUID="+uuid)
		o, err := c.CombinedOutput()
		if err != nil {
			return "", err
		}
		return string(o), err
	}
	return "", NewCommandMissingError()
}
