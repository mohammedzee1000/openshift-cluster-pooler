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

### Setup DB

Build db-cli

```
go build cmd/test-db-cli/main/db-cli.go
```

Make a test DB directory
```
mkdir `pwd`/test-badger
```

**Warning:** `db-cli` is a quick and dirty cli built to overcome missing cli for badger db.
It has limited function and should not be considered as replacement for the admin cli
which is in the pipeline

### Copy minishift provision scripts

Add minishift pool config example to db

```
BADGER_DIR="`pwd`/test-badger" ./db-cli `pwd`/openshift-cluster-pooler/pool-examples/minishift-simple/minishift-simple.json
```

### Run pool manager

Copy scripts to /usr/bin/
```
$ cp -avrf pool-examples/minishift-simple/usr/bin/* /usr/bin/
```

Build pool manager

```
$ go build cmd/pool-manager/main/pool-manager.go
```


Start pool manager

```
BADGER_DIR="`pwd`/test-badger" ./pool-manager
```
