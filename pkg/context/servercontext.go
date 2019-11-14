package context

import (
	"io/ioutil"
	"os"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//ServiceContext repersents context for services
type ServiceContext struct {
	dataDir string
	*coreContext
}

//NewServiceContext creates a new ServiceContext
func NewServiceContext(name string, dataDir string, logsDir string, debug bool) *ServiceContext {
	if len(dataDir) <= 0 {
		dataDir = "/tmp/badger"
	}
	_ = os.MkdirAll(dataDir, os.ModePerm)
	return &ServiceContext{
		dataDir:     dataDir,
		coreContext: newCoreContext(name, logsDir, debug),
	}
}

//newDataConnection creates a new data connection
func (sc *ServiceContext) newDataConnection() (*badger.DB, error) {
	bo := badger.DefaultOptions(sc.dataDir)
	if !sc.debug {
		emptyLogger := logrus.New()
		emptyLogger.Out = ioutil.Discard
		bo = bo.WithLogger(emptyLogger)
	}
	db, err := badger.Open(bo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//DataTransaction creates a new data storage transaction
func (sc *ServiceContext) DataTransaction(write bool, fn func(tx *badger.Txn) error) error {
	db, err := sc.newDataConnection()
	if err != nil {
		return errors.Wrap(err, "unable to complete connection to database")
	}
	defer db.Close()
	if !write {
		err = db.View(fn)
	} else {
		err = db.Update(fn)
	}
	if err != nil {
		return errors.Wrap(err, "unable to complete transaction")
	}
	return nil
}
