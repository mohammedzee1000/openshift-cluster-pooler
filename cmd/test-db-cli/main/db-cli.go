package main

import (
	"encoding/json"
	"io/ioutil"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/config"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/database"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/prometheus/common/log"
	"os"
)

func main()  {
	op := os.Args[1]
	ctx, err := config.NewContext()
	if err != nil {
		log.Fatal(err.Error())
	}
	switch op {
	case "loadpool":
		pool := pools.NewEmptyPool()
		fl := os.Args[2]
		data, err := ioutil.ReadFile(fl)
		if err != nil {
			log.Fatalf("unable to read data %s", err.Error())
		}

		err = json.Unmarshal(data, &pool)
		if err != nil {
			log.Fatalf("unable to unmarshal pool %s", err.Error())
		}
		key := pools.GetPoolKey(pool.Name)
		database.SaveinKVDB(&ctx, key, string(data))
		break
	case "getpool":
		break
	default:
		log.Fatal("invalid option, try `loadpool filename` or getpool poolname")

	}
}
