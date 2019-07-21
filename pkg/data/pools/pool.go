package pools

import (
	"errors"
	"sync"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/config"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/logging"
	"time"
)

func calculateTimeoutPeriod(period int) time.Duration {
	tp := time.Hour
	for i := 1; i < period; i++ {
		tp = tp + time.Hour
	}
	return tp
}

//GC deprovisions clusters that have failed or outlived a timeout
func (p Pool) GC(ctx *config.Context) error  {
	logging.Info("Pool GC", "initiating GC of pool %s", p.Name)
	//Fetch all clusterlist
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
	if len(gcclusters) <= 0 {
		logging.Info("Pool GC", "no garbage in pool %s, skipping", p.Name)
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
						logging.Error("Pool GC", "Failed to force deprovision cluster : %s", err.Error())
					}
				} else {
					logging.Error("Pool GC", "Failed to deprovision cluster : %s", err.Error())
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

//Reconcole ensured that expected and actual pool size match
func (p Pool) Reconcile(ctx *config.Context) error {
	logging.Info("Pool Reconcile", "initiating reconciliation for pool %s", p.Name)
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
		logging.Info("Pool Reconcile", "available clusters do not match expected for pool %s", p.Name)
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
			logging.Info("Pool Reconcile", "allocating %d clusters for pool %s serially", tp, p.Name)
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
					return errors.New("failed to provision some clusters")
					// todo use this for some sore of back off login
				}
			}
		}
	} else {
		logging.Info("Pool Reconcile", "skipping reconcilation for pool %s as actual matches expected", p.Name)
	}
	return nil
}