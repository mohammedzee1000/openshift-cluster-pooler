package config

import (
	"context"
	"github.com/dgraph-io/badger"
	"os"
	"time"
)

type Context struct {
	*internalContext
}

type internalContext struct {
	LogsDir       string
	BadgerDir     string
}

func NewCliContext() context.Context {
	c, _ := context.WithTimeout(context.Background(), 15*time.Second)
	return c
}

func NewContext() (Context, error) {
	ctx := Context{&internalContext{}}
	ctx.LogsDir = "/var/log/openshift-clusters-pools"
	ctx.BadgerDir = "/var/openshift-cluster-pools/badger"
	badgerenv := os.Getenv("BADGER_DIR")
	if len(badgerenv) > 0 {
		ctx.BadgerDir = badgerenv
	}
	return ctx, nil
}


func (c Context) NewBadgerConnection() (*badger.DB, error)  {
	db, err := badger.Open(badger.DefaultOptions(c.BadgerDir))
	if err != nil {
		return nil, err
	}
	return db, nil
}