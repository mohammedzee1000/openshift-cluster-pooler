package clusters

import (
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"time"

)

const (
	State_Provisioning   = "Provisioning"
	State_Failed         = "Failed"
	State_Success        = "Success"
	State_Used           = "Activated"
	State_DeProvisioning = "Deprovisioning"
	State_Cleanup        = "Cleanup"
	State_Returned       = "Returned"
	Cluster_Prefix       = "Cluster"
)

type Cluster struct {
	ClusterID   string `json:"ClusterID"`
	PoolName    string `json:"PoolName"`
	State       string `json:"State"`
	URL  string        `json:"URL"`
	AdminUser   string `json:"AdminUser"`
	AdminPassword string `json:"AdminPassword"`
	CAFile 		[]string `json:"CAFile"`
	CertFile    []string `json:"CertFile"`
	KeyFile 	[]string `json:"KeyFile"`
	ExtraInfo	string `json:"ExtraInfo"`
	CreatedOn   time.Time `json:"CreatedOn"`
	ActivatedOn time.Time `json:"ActivatedOn"`
}

func NewEmptyCluster() *Cluster {
	return &Cluster{}
}

func NewCluster(clusterid string, poolName string) *Cluster {
	c := NewEmptyCluster()
	c.ClusterID = clusterid
	c.PoolName = poolName
	c.State = State_Provisioning
	return c
}

func Equal(x *Cluster, y *Cluster) bool {
	if x.ClusterID == y.ClusterID {
		return true
	}
	return false
}

func DeepEqual(x *Cluster, y *Cluster) bool {
	//Todo implement
	return false
}

type ClusterList struct {
	items []*Cluster
}

func NewClusterList() *ClusterList {
	var clusters []*Cluster
	return WrapClusterList(clusters)
}

func WrapClusterList(clusters []*Cluster) *ClusterList  {
	return &ClusterList{items: clusters}
}

func (c *Cluster) Return(ctx *generic.Context) error  {
	c.State = State_Returned
	return c.Save(ctx)
}