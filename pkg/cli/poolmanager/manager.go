package poolmanager

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/config"
	pools2 "github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/logging"
	"github.com/prometheus/common/log"
	"time"
)

type PoolManager struct {
	context *config.Context
}

func NewPoolManager() *PoolManager {
	ctx, err := config.NewContext()
	if err != nil {
		log.Fatal("unable to initialize context")
	}
	return &PoolManager{context: &ctx}
}

func (pm PoolManager) Run()  {
	//Initialize
	c, err := pm.context.NewBadgerConnection()

	if err != nil {
		log.Fatal("unable to connect to database : ", err.Error())
	}
	_ = c.Close()
	for {
		pools, err := pools2.List(pm.context)
		if err != nil {
			logging.Fatal("Pool Manager","failed to retrieve pool list : ", err.Error())
		}
		if len(pools) > 0 {
			for _, item := range pools{
				_ = item.GC(pm.context)
				_ = item.Reconcile(pm.context)
				time.Sleep(1 * time.Minute)
			}
		} else {
			logging.Info("Pool Manager", "no pools, skipping this turn")
		}
		time.Sleep(2 * time.Minute)
	}
}