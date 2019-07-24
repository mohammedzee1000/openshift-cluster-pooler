package pools

import (
	"errors"
	"fmt"
	"sync"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"time"
)


func (p Pool) gcCollect(ctx *generic.Context, componentSubName string,gcclusters []clusters.Cluster) error  {
	var err error
	componentName := fmt.Sprintf("Pool GC - %s", componentSubName)
	if len(gcclusters) <= 0 {
		ctx.Log.Info(componentName, "no garbage in pool %s, skipping", p.Name)
		return nil
	}
	// Delete clusterlist as needed
	// If no paralled deprovisioning, do it serially
	if p.ParallelDeProvisioning <= 1 {
		for _, c := range gcclusters {
			err = p.deprovision(ctx, c.ClusterID, false)
			if err != nil {
				if len(p.ForceDeprovisionCommand) > 0 {
					err = p.deprovision(ctx, c.ClusterID, true)
					if err != nil {
						ctx.Log.Error(componentName, err, "Failed to force deprovision cluster")
					}
				} else {
					ctx.Log.Error(componentName, err, "Failed to deprovision cluster")
				}
			}
		}
	} else {
		//Example paralleldeprovision = 3, total = 5
		//first iter - todeprovision = 3, total = 2
		//second iter - todeprovision = 3, total = -1; todeprovision = 3-1 = 2
		total := len(gcclusters)
		for total > 0 {
			//assume we need to deprovision, parallel deprovision no of times
			todeprovision := p.ParallelDeProvisioning
			//remove it from total
			total = total - todeprovision
			//if total is negative, then we will overprovision to remove total from provision
			if total < 0 {
				todeprovision = todeprovision + total
			}
			//deprovision in parallel
			chanerrors := make(chan error, todeprovision)
			wg := new(sync.WaitGroup)
			wg.Add(todeprovision-1)
			for i := 0; i < todeprovision; i++ {
				go func() {
					err := p.deprovision(ctx, gcclusters[i].ClusterID, false)
					if err != nil {
						if len(p.ForceDeprovisionCommand) > 0 {
							err = p.deprovision(ctx, gcclusters[i].ClusterID, true)
							if err != nil {
								chanerrors <- err
							}
						} else {
							chanerrors <- err
						}
					}
					wg.Done()
				}()
			}
			wg.Wait()
			if len(chanerrors) > 0 {
				return errors.New("failed to deprovision some clusters")
				// todo this length should be used in backoff
			}
			close(chanerrors)
		}
	}
	return nil
}

func (p Pool) gcByCondition(ctx *generic.Context) error  {
	var gcclusters []clusters.Cluster
	clusterlist, err := clusters.ClustersInPool(ctx, p.Name)
	if err != nil {
		return err
	}
	// gather clusterlist to delete
	for _, c := range clusterlist {
		if c.State == clusters.State_Failed || c.State == clusters.State_Cleanup {
			gcclusters = append(gcclusters, c)
		} else if c.State == clusters.State_Used || c.State == clusters.State_Success {
			//dead time, knows when the clusterlist is supposed to have died
			var dt time.Time
			// Calculate time passed as needed
			if c.State == clusters.State_Success {
				dt = c.CreatedOn.Add(time.Hour * time.Duration(p.UnusedClusterTimeout))
			} else {
				dt = c.ActivatedOn.Add(time.Hour * time.Duration(p.UsedClusterTimeout))
			}

			//if current time is after dead time, cleanup !!
			if dt.Before(time.Now()) {
				gcclusters = append(gcclusters, c)
			}
		}
	}
	return p.gcCollect(ctx, "By Condition",gcclusters)
}

func (p Pool) gcByConfigChange(ctx *generic.Context) error {
	clusterlist, err := clusters.ClustersInPool(ctx, p.Name)
	if err != nil {
		return err
	}
	// find out how many clusters need to be removed
	// len of success clusters and used clusters
	successclusters := clusters.ClustersInStateIn(clusterlist, clusters.State_Success)
	successclustercount := len(successclusters)
	usedclusterscount := len(clusters.ClustersInStateIn(clusterlist, clusters.State_Used))
	var toremove int
	if p.MaxSize > p.Size && usedclusterscount >= p.Size - 1 {
		toremove = successclustercount - p.MaxSize
	}  else {
		toremove = successclustercount - p.Size
	}
	return p.gcCollect(ctx, "By Config Change", clusters.OldestN(clusterlist, toremove))
}

//GC deprovisions clusters that have failed or outlived a timeout
func (p Pool) GC(ctx *generic.Context) error  {
	ctx.Log.Info("Pool GC", "initiating GC of pool %s", p.Name)
	//Fetch all clusterlist
	ctx.Log.Info("Pool GC", "initiating cleanup of clusters that have met some conditions, pool %s", p.Name)
	_ = p.gcByCondition(ctx)
	ctx.Log.Info("Pool GC", "initiating cleanup of clusters that need to be removed due to config change, pool %s", p.Name)
	_ = p.gcByConfigChange(ctx)
	return nil
}


//Reconcile ensured that expected and actual pool size match
func (p Pool) Reconcile(ctx *generic.Context) error {
	ctx.Log.Info("Pool Reconcile", "initiating reconciliation for pool %s", p.Name)
	//Fetch all clusterlist
	currentClusters := 0
	activatedClusters := 0
	clusterlist, err := clusters.ClustersInPool(ctx, p.Name)
	if err != nil {
		return err
	}
	//Calculate current clusterlist count
	for _, c := range clusterlist {
		if c.State == clusters.State_Provisioning || c.State == clusters.State_Success || c.State == clusters.State_Used {
			currentClusters = currentClusters + 1
			if c.State == clusters.State_Used {
				activatedClusters = activatedClusters + 1
			}
		}
	}
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