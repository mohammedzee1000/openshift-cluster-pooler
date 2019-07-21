# openshift-cluster-pool
Code base for openshift cluster pool

## Try It (With minishift) - Developer

### Prepare

1. Make sure you have latest golang, docker and etcdctl the cli for etcd installed

1. Make sure you have minishift installed (https://docs.okd.io/latest/minishift/getting-started/index.html)

1. Clone repo
```
$ git clone mohammedzee1000/openshift-cluster-pool && cd openshift-cluster-pool
```

### Setup and run

Copy scripts to /usr/bin/
```
$ cp -avrf pool-examples/minishift-simple/usr/bin/* /usr/bin/
```

Start etcd server

```
$ sudo docker run -d -p 2379:2379 -p 4001:4001 quay.io/coreos/etcd:latest etcd --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001
```

Build pool manager

```
$ go build cmd/pool-manager/main/pool-manager.go
```

Load minishift pool to db

```
$ echo pool-examples/minishift-simple/minishift-simple.json | etcdctl put Pool-minishift-simple
```

Start pool manager

```
$ ETCD_ENDPOINTS="http://0.0.0.0:2379" ./pool-manager
```
