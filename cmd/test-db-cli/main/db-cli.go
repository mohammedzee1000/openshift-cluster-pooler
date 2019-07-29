package main

import (
	"encoding/json"
	"fmt"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
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
	case "save-pool":
		pool := pools.NewEmptyPool()
		fl := os.Args[2]
		data, err := ioutil.ReadFile(fl)
		if err != nil {
			ctx.Log.Fatal("save-pool", err, "unable to read data")
		}

		err = json.Unmarshal(data, &pool)
		if err != nil {
			ctx.Log.Fatal("save-pool", err, "unable to unmarshal pool")
		}
		key := pools.GetPoolKey(pool.Name)
		database.SaveinKVDB(ctx, key, string(data))
		ctx.Log.Info("save-pool", "Loading pool into DB")
		break
	case "get-pool":
		pnm := os.Args[2]
		p, err := pools.PoolByName(ctx, pnm, false)
		if err != nil {
			ctx.Log.Fatal("get-pool", err,"could not get pool info")
		}
		d, err := json.Marshal(p)
		if err != nil {
			log.Fatal("get-pool", err, "failed to unmarshal data")
		}
		fmt.Println(string(d))
		break
	case "get-pool-del":
		pnm := os.Args[2]
		p, err := pools.PoolByName(ctx, pnm, true)
		if err != nil {
			ctx.Log.Fatal("get-pool-del", err,"could not get pool info")
		}
		d, err := json.Marshal(p)
		if err != nil {
			log.Fatal("get-pool-del", err, "failed to unmarshal data")
		}
		fmt.Println(string(d))
		break
	case "del-pool":
		pnd := os.Args[2]
		p, err := pools.PoolByName(ctx, pnd, false)
		if err != nil{
			ctx.Log.Fatal("del-pool", err, "failed to retrieve pool of that name")
		}
		err = p.MarkForRemoval(ctx)
		if err != nil {
			ctx.Log.Fatal("del-pool", err, "failed to mark pool for removal")
		}
		ctx.Log.Info("del-pool", "marked pool %s for removal. it will be cleaned in next gc cycle by pool manager", pnd)
	case "list-clusters":
		cl, err := clusters.List(ctx)
		if err != nil {
			ctx.Log.Fatal("list-clusters", err, "cannot list clusters")
		}
		for _, item := range cl.Items{
			data, err := json.Marshal(item)
			if err != nil {
				ctx.Log.Error("list-clusters", err, "failed to unmarshal cluster info")
			}
			fmt.Println(data)
		}
	default:
		sp := "save-pool <path>"
		gp := "get-pool <name>"
		gpd := "get-pool-del <name>"
		dp := "del-pool <name>"
		lc := "list-clusters"
		fmt.Printf("invalid command, following are possible:\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s", sp, gp, gpd, dp, lc)
	}
}
