# openshift-cluster-pool
Code base for openshift cluster pool

## Try It (With minishift) - Developer

### Prepare

1. Make sure you have latest golang setup with support for go mod.

1. Make sure you have minishift installed (https://docs.okd.io/latest/minishift/getting-started/index.html)

1. Clone repo
```
$ git clone mohammedzee1000/openshift-cluster-pool && cd openshift-cluster-pool
```

**Note:** You may need to setup requirements in go.mod appropriately before proceeding

### Build the binaries

Build all the nessasary binaries with 

```bash
$ make build
# or
# scripts/build.sh
```

### Setup DB

Make a test DB directory
```
mkdir `pwd`/test-badger
```
Load sample pool config

```
BADGER_DIR="`pwd`/test-badger" ./db-cli save-pool `pwd`/pool-examples/minishift-simple/minishift-simple.json
```

**Warning:** `db-cli` is a quick and dirty cli built to overcome missing cli for badger db.
It has limited function and should not be considered as replacement for the admin cli
which is in the pipeline

### Copy minishift provision scripts

Copy scripts to /usr/bin/
```
$ cp -avrf pool-examples/minishift-simple/usr/bin/* /usr/bin/
```

### Start pool manager

```bash
BADGER_DIR="`pwd`/test-badger" ./pool-manager
```

### Other important options of db-cli

As already stated `db-cli` is only for testing purposes. But you can use still use it.
Just do `./db-cli help` to find other commands

### Client API Server

A basic client api server is provided that can be launched with

```bash
$ BADGER_DIR="`pwd`/test-badger" HOST_ON=":20000" ./api-server
```

#### Possible operations:

##### List pool
 
 ```bash
$ curl http://localhost:20000/pools/list
```

##### Describe a pool

```bash
$ curl http://localhost:20000/pool/minishift-simple/describe
```

##### Activate cluster

```bash
$ curl http://localhost:20000/pool/minishift-simple/get-cluster
```

##### Describe Cluster

```bash
$ curl -k http://localhost:20000/cluster/{clusterid}/describe
```

##### Return used cluster early for cleanup

```bash
$ curl -k http://localhost:20000/cluster/{clusterid}/return
```

## In the pipeline
 - Admin API Server + client to replace db-cli
 - Client API server + client to allow users to list pools and get a cluster from pool / list its information
