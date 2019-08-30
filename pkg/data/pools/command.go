package pools

import (
	"errors"
	"fmt"
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

func PrintIfDebug(debug bool, title string, out string)  {
	if debug {
		fmt.Printf("%s : \n%s\n", title, out)
	}
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
		return strings.TrimRight(string(o), "\n"), err
	}
	return "", NewCommandMissingError()
}
