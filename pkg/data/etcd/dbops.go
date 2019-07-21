package etcd

import (
	"context"
	"github.com/etcd-io/etcd/etcdserver/api/v3rpc/rpctypes"
	"log"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/config"
)

// TODO make it retry for etcd stuff

func handleError(err error)  {
		if err != nil {
			switch err {
			case context.Canceled:
				log.Fatalf("ctx is canceled by another routine: %v", err)
			case context.DeadlineExceeded:
				log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
			case rpctypes.ErrEmptyKey:
				log.Fatalf("client-side error: %v", err)
			default:
				log.Fatalf("bad clusters endpoints, which are not etcd servers: %v", err)
			}
		}
}

//SaveinEtcd saved specified key value pair in etcd
func SaveInEtcd(ctx *config.Context, key string, data string) {
	cli, err := ctx.NewEtcdConnection()
	if err != nil {
		handleError(err)
	}
	defer cli.Close()
	_, err = clientv3.NewKV(cli).Put(config.NewCliContext(), key, data)
	if err != nil {
		handleError(err)
	}
}

//GetMultipleWithPrefixFromEtcd gets multiple values whose keys match specified prefix in etcd
func GetMultipleWithPrefixFromEtcd(ctx *config.Context, keyprefix string) []string {
	var values []string
	cli, err := ctx.NewEtcdConnection()
	if err != nil {
		handleError(err)
	}
	defer cli.Close()
	resp, err := clientv3.NewKV(cli).Get(config.NewCliContext(), keyprefix, clientv3.WithPrefix())
	if err != nil {
		handleError(err)
	}
	for _, item := range resp.Kvs {
		values = append(values, string(item.Value))
	}
	return values
}

//GetExactFromEtcd gets specific value which matches exact string
func GetExactFromEtcd(ctx *config.Context, key string) string {
	var value string
	cli, err := ctx.NewEtcdConnection()
	if err != nil {
		handleError(err)
	}
	defer cli.Close()
	resp, err := clientv3.NewKV(cli).Get(config.NewCliContext(), key)
	if err != nil {
		handleError(err)
	}
	if len(resp.Kvs) > 0 {
		value = string(resp.Kvs[0].Value)
	}
	return value
}

//DeleteInEtcd deletes the key specified in etcd
func DeleteInEtcd(ctx *config.Context, key string)  {
	cli, err := ctx.NewEtcdConnection()
	if err != nil {
		handleError(err)
	}
	defer cli.Close()
	_, err = clientv3.NewKV(cli).Delete(config.NewCliContext(), key)
	if err != nil {
		handleError(err)
	}
}