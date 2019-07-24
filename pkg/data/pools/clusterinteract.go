package pools

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
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
