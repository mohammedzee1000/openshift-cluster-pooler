package config

import (
	"context"
	"github.com/etcd-io/etcd/clientv3"
	"os"
	"strings"
	"time"
)

type Context struct {
	*internalContext
}

type internalContext struct {
	LogsDir       string
	EtcdEndpoints []string
}

func NewCliContext() context.Context {
	c, _ := context.WithTimeout(context.Background(), 15*time.Second)
	return c
}

func NewContext() (Context, error) {
	ctx := Context{&internalContext{}}
	ctx.LogsDir = "/var/log/openshift-clusters-pools"
	etcdEndPointsEnv := os.Getenv("ETCD_ENDPOINTS")
	if len(etcdEndPointsEnv) > 0 {
		ctx.EtcdEndpoints = strings.Split(etcdEndPointsEnv, ";")
	}
	return ctx, nil
}

func (c Context) NewEtcdConnection() (*clientv3.Client ,error) {
	cl, err := clientv3.New(clientv3.Config{
		Endpoints: c.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cl, nil
}


// template to be removed
//func (c Context) Get(key string) (string, error) {
//	cl, err := c.NewEtcdConnection()
//	if err != nil {
//		return "", err
//	}
//	defer cl.Close()
//	kv := clientv3.NewKV(cl)
//	resp, err := kv.Get(NewCliContext(), key)
//	if err != nil {
//		switch err {
//		case context.Canceled:
//			log.Fatalf("ctx is canceled by another routine: %v", err)
//		case context.DeadlineExceeded:
//			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
//		case rpctypes.ErrEmptyKey:
//			log.Fatalf("client-side error: %v", err)
//		default:
//			log.Fatalf("bad clusters endpoints, which are not database servers: %v", err)
//		}
//		return "", err
//	}
//
//	return string(resp.Kvs[0].Value), nil
//}
