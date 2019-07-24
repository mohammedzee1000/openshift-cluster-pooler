package main

import (
	"encoding/json"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/database"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"os"
)

func main()  {
	if len(os.Args) < 1{
		log.Fatalln("please pass a parameter")
	}
	op := os.Args[1]
	ctx, err := generic.NewContext("test-db-cli")
	if err != nil {
		log.Fatal(err.Error())
	}
	//this if is just to remove warning
	if ctx == nil {
		return
	}

	switch op {
	case "loadpool":
		pool := pools.NewEmptyPool()
		fl := os.Args[2]
		data, err := ioutil.ReadFile(fl)
		if err != nil {
			ctx.Log.Fatal("loadpool", err, "unable to read data")
		}

		err = json.Unmarshal(data, &pool)
		if err != nil {
			ctx.Log.Fatal("loadpool", err, "unable to unmarshal pool")
		}
		key := pools.GetPoolKey(pool.Name)
		database.SaveinKVDB(ctx, key, string(data))
		break
	case "getpool":
		break
	default:
		log.Fatal("invalid option, try `loadpool filename` or getpool poolname")

	}
}
