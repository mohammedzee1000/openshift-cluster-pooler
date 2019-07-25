package pools

import (
	"fmt"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"github.com/pkg/errors"
	"sync"
	"time"
)

func (p Pool) gcCollect(ctx *generic.Context, componentSubName string,gcclusters *clusters.ClusterList) error  {
	var err error
	componentName := fmt.Sprintf("Pool GC - %s", componentSubName)
	if len(gcclusters.Items) <= 0 {
		ctx.Log.Info(componentName, "no garbage in pool %s, skipping", p.Name)
		return nil
	}
	// deleteInDB clusterlist as needed
	// If no paralled deprovisioning, do it serially
	if p.ParallelDeProvisioning <= 1 {
		for _, item := range gcclusters.Items {
			err = p.deprovision(ctx, item.ClusterID, false)
			if err != nil {
				if len(p.ForceDeprovisionCommand) > 0 {
					err = p.deprovision(ctx, item.ClusterID, true)
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
		total := len(gcclusters.Items)
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
					err := p.deprovision(ctx, gcclusters.Items[i].ClusterID, false)
					if err != nil {
						if len(p.ForceDeprovisionCommand) > 0 {
							err = p.deprovision(ctx, gcclusters.Items[i].ClusterID, true)
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
	gcclusters := clusters.NewClusterList()
	clusterlist, err := p.GetClusters(ctx)
	if err != nil {
		return err
	}
	// gather clusterlist to delete
	for _, item := range clusterlist.Items {
		if item.State == clusters.State_Failed || item.State == clusters.State_Cleanup {
			gcclusters.Items = append(gcclusters.Items, item)
		} else if item.State == clusters.State_Used || item.State == clusters.State_Success {
			//dead time, knows when the clusterlist is supposed to have died
			dt := p.ExpiresOn(&item)
			//if current time is after dead time, cleanup !!
			if dt.Before(time.Now()) {
				gcclusters.Items = append(gcclusters.Items, item)
			}
		}
	}
	return p.gcCollect(ctx, "By Condition", gcclusters)
}

func (p Pool) gcByConfigChange(ctx *generic.Context) error {
	clusterlist, err := p.GetClusters(ctx)
	if err != nil {
		return err
	}
	// find out how many clusters need to be removed
	// len of success clusters and used clusters
	successclusters := clusterlist.ClustersInStateIn(clusters.State_Success)
	successclustercount := len(successclusters.Items)
	usedclusterscount := len(clusterlist.ClustersInStateIn(clusters.State_Used).Items)
	var toremove int
	if p.MaxSize > p.Size && usedclusterscount >= p.Size - 1 {
		toremove = successclustercount - p.MaxSize
	}  else {
		toremove = successclustercount - p.Size
	}
	return p.gcCollect(ctx, "By Config Change", clusterlist.OldestN(toremove))
}

//Destroys the clusters and the pool
func (p Pool) Destroy(ctx *generic.Context) error {
	if p.IsMarkedForRemoval(ctx) {
		clusterlist, err := p.GetClusters(ctx)
		if err != nil {
			ctx.Log.Error("Pool Remover", err, "failed to retrieve clusters in pool %s", p.Name)
		}
		ctx.Log.Info("Pool Remover", "initiating cleanup of clusters in pool %s", p.Name)
		err = p.gcCollect(ctx, "Pool Remover", clusterlist)
		if err != nil {
			ctx.Log.Error("Pool Remover", err, "failed to destroy pool %s", p.Name)
		}
		p.deleteInDB(ctx)
		return nil
	}
	e := errors.New("this pool is not marked for removal")
	ctx.Log.Error("Pool Remover", e, "")
	return e
}

//GC deprovisions clusters that have failed or outlived a timeout
func (p Pool) GC(ctx *generic.Context) error  {
	ctx.Log.Info("Pool GC", "initiating GC of pool %s", p.Name)
	ctx.Log.Info("Pool GC", "initiating cleanup of clusters that have met some conditions, pool %s", p.Name)
	_ = p.gcByCondition(ctx)
	ctx.Log.Info("Pool GC", "initiating cleanup of clusters that need to be removed due to config change, pool %s", p.Name)
	_ = p.gcByConfigChange(ctx)
	return nil
}