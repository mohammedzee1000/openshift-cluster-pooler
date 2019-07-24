package poolmanager

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	pools2 "github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"time"
	"log"
)

type PoolManager struct {
	context *generic.Context
}

func NewPoolManager() *PoolManager {
	ctx, err := generic.NewContext("pool-manager")
	if err != nil {
		log.Fatal("Pool Manager", err, "unable to initialize context")
	}
	return &PoolManager{context: ctx}
}

func (pm PoolManager) Run()  {
	//Initialize
	c, err := pm.context.NewBadgerConnection()

	if err != nil {
		pm.context.Log.Fatal("Pool Manager", err, "unable to connect to database")
	}
	_ = c.Close()
	for {
		pools, err := pools2.List(pm.context)
		if err != nil {
			pm.context.Log.Fatal("Pool Manager", err, "failed to retrieve pool list")
		}
		//todo have a way to gc clusters whose config is removed
		if len(pools) > 0 {
			for _, item := range pools{
				_ = item.GC(pm.context)
				_ = item.Reconcile(pm.context)
				time.Sleep(1 * time.Minute)
			}
		} else {
			pm.context.Log.Info("Pool Manager", "no pools, skipping this turn")
		}
		time.Sleep(2 * time.Minute)
	}
}