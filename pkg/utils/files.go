package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/util"
	"github.com/pkg/errors"
)

//WriteFile writes data to a file
func WriteFile(filepath string, data []byte) error {
	return ioutil.WriteFile(filepath, data, os.ModePerm)
}

//ReadYamlFile reads the contents of yaml file and returns the contents
func ReadYamlFile(o interface{}, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "unable to read yaml file")
	}
	return yaml.Unmarshal(data, o)
}

//EvalStringWithUUID uses shell evaluation to evaluate an output
func EvalStringWithUUID(uuid string, str string) (string, error) {
	if len(str) <= 0 {
		return "", nil
	}
	command := fmt.Sprintf("echo %s", str)
	return RunCommandWithUUID(uuid, command)
}

//GetFileContent reads a file content while evaluating its output
func GetFileContent(clusterid string, path string) ([]string, error) {
	out, err := EvalStringWithUUID(clusterid, path)
	if err != nil {
		return nil, err
	}
	return util.ReadFileLines(string(out))
}
