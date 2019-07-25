package poolmanager

import (
	pools2 "github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"log"
	"time"
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
	pm.context.Log.Info("Pool Manager", "starting...")
	for {
		pm.context.Log.Info("Pool Manager", "initiating cluster management cycle")
		pm.context.Log.Info("Pool Manager", "managing pools to be removed")
		rpools, err := pools2.List(pm.context, true)
		if err != nil {
			pm.context.Log.Error("Pool Manager", err, "failed to retrieve gc pool list")
		}
		if len(rpools) > 0 {
			for _, item := range rpools {
				_ = item.Destroy(pm.context)
			}
		} else {
			pm.context.Log.Info("Pool Manager", "no pools to remove, skipping this turn")
		}
		pm.context.Log.Info("Pool Manager", "managing available pools")
		pools, err := pools2.List(pm.context, false)
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