package main

import (
	"encoding/json"
	"fmt"
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
		ctx.Log.Info("loadpool", "Loading pool into DB")
		break
	case "getpool":
		pnm := os.Args[2]
		p, err := pools.PoolByName(ctx, pnm, false)
		if err != nil {
			ctx.Log.Fatal("getpool", err,"could not get pool info")
		}
		d, err := json.Marshal(p)
		if err != nil {
			log.Fatal("getpool", err, "failed to unmarshal data")
		}
		fmt.Println(string(d))
		break
	case "getpooldel":
		pnm := os.Args[2]
		p, err := pools.PoolByName(ctx, pnm, true)
		if err != nil {
			ctx.Log.Fatal("getpooldel", err,"could not get pool info")
		}
		d, err := json.Marshal(p)
		if err != nil {
			log.Fatal("getpooldel", err, "failed to unmarshal data")
		}
		fmt.Println(string(d))
		break
	case "delpool":
		pnd := os.Args[2]
		p, err := pools.PoolByName(ctx, pnd, false)
		if err != nil{
			ctx.Log.Fatal("delpool", err, "failed to retrieve pool of that name")
		}
		err = p.MarkForRemoval(ctx)
		if err != nil {
			ctx.Log.Fatal("delpool", err, "failed to mark pool for removal")
		}
		ctx.Log.Info("delpool", "marked pool %s for removal. it will be cleaned in next gc cycle by pool manager", pnd)
	default:
		log.Fatal("invalid option, try `loadpool filename` or getpool poolname")

	}
}
