package generic

import (
	"context"
	"github.com/dgraph-io/badger"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

type Context struct {
	*internalContext
}

type internalContext struct {
	Name          string
	LogsDir       string
	BadgerDir     string
	Debug         bool
	Log           *LogHander
}

func NewCliContext() context.Context {
	c, _ := context.WithTimeout(context.Background(), 15*time.Second)
	return c
}

func NewContext(name string) (*Context, error) {
	ctx := Context{&internalContext{}}
	ctx.Name = name
	ctx.LogsDir = "/var/log/openshift-clusters-pools"
	ctx.BadgerDir = "/tmp/badger"
	badgerenv := os.Getenv("BADGER_DIR")
	debugenv := os.Getenv("DEBUG")
	if len(badgerenv) > 0 {
		ctx.BadgerDir = badgerenv
	}
	if len(debugenv) > 0 && debugenv == "true" {
		ctx.Debug = true
	}
	ctx.Log = NewLogger(name, ctx.Debug)
	return &ctx, nil
}


func (c Context) NewBadgerConnection() (*badger.DB, error)  {
	// todo make it log to a file instead
	bo := badger.DefaultOptions(c.BadgerDir)
	if !c.Debug {
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