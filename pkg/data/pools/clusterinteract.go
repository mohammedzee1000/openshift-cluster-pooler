package pools

import (
	"encoding/json"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/database"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"time"
)

//Gets when a cluster will expire in conjunction with a pool config
func (p Pool) ExpiresOn(c *clusters.Cluster) time.Time {
	var t time.Time
	if c.State == clusters.State_Success {
		t = c.CreatedOn.Add(time.Duration(p.UnusedClusterTimeout) * time.Hour)
	} else if c.State == clusters.State_Used {
		t = c.ActivatedOn.Add(time.Duration(p.UsedClusterTimeout) * time.Hour)
	}
	return t
}

func (p Pool) GetClusters(ctx *generic.Context) (*clusters.ClusterList, error)  {
	clusterlist := clusters.NewClusterList()
	d := database.GetMultipleWithPrefixFromKVDB(ctx, clusters.GetClusterPoolKey(p.Name))
	for _, item := range d {
		var cl clusters.Cluster
		err := json.Unmarshal([]byte(item), &cl)
		if err != nil {
			return nil, err
		}
		clusterlist.Append(&cl)
	}
	return clusterlist, nil
}