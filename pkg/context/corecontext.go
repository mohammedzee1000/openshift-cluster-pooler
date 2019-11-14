package context

import (
	"fmt"
	"os"
	"path/filepath"

	"openshift-cluster-pooler/pkg/logging"
	"openshift-cluster-pooler/pkg/utils"
)

type coreContext struct {
	logsDir string
	debug   bool
	logger  *logging.LogHandler
}

func newCoreContext(name string, logsDir string, debug bool) *coreContext {
	return &coreContext{
		logsDir: logsDir,
		debug:   debug,
		logger:  logging.NewLogger(name, debug),
	}
}

func (c coreContext) LogToFile(filename string, data []byte, subdirs ...string) error {
	var dl string
	dl = filepath.Join(c.logsDir, subdirs[0])
	for _, elem := range subdirs[1:] {
		dl = filepath.Join(dl, elem)
	}
	fl := filepath.Join(dl, filename)
	_ = os.MkdirAll(dl, os.ModePerm)
	c.logger.Info("file logger", "writing logs to file %s", fl)
	if c.debug {
		fmt.Println("logs : \n", string(data))
	}
	return utils.WriteFile(fl, data)
}

func (c coreContext) DebugEnabled() bool {
	return c.debug
}

func (c coreContext) GetLogger() *logging.LogHandler {
	return c.logger
}
