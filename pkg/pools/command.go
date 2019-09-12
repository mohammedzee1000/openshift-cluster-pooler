package pools

import (
	"errors"
	"fmt"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/util"
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

func getFileContent(clusterid string, path string) ([]string, error)  {
	out, err := evalString(clusterid, path)
	if err != nil {
		return nil, err
	}
	return util.ReadFileLines(string(out))
}

func evalString(uuid string, str string) (string, error)  {
	if len(str) <= 0 {
		return "", nil
	}
	command := fmt.Sprintf("echo %s", str)
	return runCommand(uuid, command)
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
