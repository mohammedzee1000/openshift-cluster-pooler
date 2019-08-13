package pools

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"github.com/pkg/errors"
	"sync"
)

//Reconcile ensured that expected and actual pool size match
func (p Pool) Reconcile(ctx *generic.Context) error {
	ctx.Log.Info("Pool Reconcile", "initiating reconciliation for pool %s", p.Name)
	//Fetch all clusterlist
	currentClusters := 0
	activatedClusters := 0
	clusterlist, err := p.GetClusters(ctx)
	if err != nil {
		return err
	}
	//Calculate current clusterlist count
	clusterlist.List(func(c *clusters.Cluster) {
		if c.State == clusters.State_Provisioning || c.State == clusters.State_Success || c.State == clusters.State_Used {
			currentClusters = currentClusters + 1
			if c.State == clusters.State_Used {
				activatedClusters = activatedClusters + 1
			}
		}
	})
	//Only if current clusterlist are less than expected
	if currentClusters < p.Size {
		ctx.Log.Info("Pool Reconcile", "available clusters do not match expected for pool %s", p.Name)
		tp := 0
		//if max pool size <= pool size then we dont even need to look at pool size
		if p.MaxSize <= p.Size {
			tp = p.Size - currentClusters
		} else {
			if(p.Size - activatedClusters == 1) {
				tp = p.MaxSize - currentClusters
			} else {
				tp = p.Size - currentClusters
			}
		}
		//provision tp clusterlist
		//if no parallel provisioning, provision in series
		if p.ParallelProvisioning <= 1 {
			ctx.Log.Info("Pool Reconcile", "allocating %d clusters for pool %s serially", tp, p.Name)
			for i:= 0; i < tp; i++ {
				err = p.provision(ctx)
				if err != nil {
					return err
				}
			}
		} else {
			//see explain in GC for similar parallel deprovisioning
			for tp > 0 {
				toprovision := p.ParallelProvisioning
				tp = tp - toprovision
				if tp < 0 {
					toprovision = toprovision + tp
				}
				//provision in parallel
				chanerror := make(chan error, tp)
				wg := new(sync.WaitGroup)
				wg.Add(toprovision - 1)
				for i:=0 ; i<toprovision ; i++ {
					go func() {
						err := p.provision(ctx)
						if err != nil {
							chanerror <- err
						}
						wg.Done()
					}()
				}
				wg.Wait()
				if len(chanerror) > 0 {
					ctx.Log.Info("Pool reconcile", "failed to provision some clusers")
					return errors.New("failed to provision some clusters")
					// todo use this for some sore of back off login
				}
			}
		}
	} else {
		ctx.Log.Info("Pool Reconcile", "skipping reconcilation for pool %s as actual matches expected", p.Name)
	}
	return nil
}

