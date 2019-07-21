package clusters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/config"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/etcd"
)

func getClusterPoolKey(poolname string) string {
	return fmt.Sprintf("%s-%s", Cluster_Prefix, poolname)
}

func getClusterKey(clusterid string, poolname string) string {
	return fmt.Sprintf("%s-%s-%s", Cluster_Prefix, poolname, clusterid)
}

//Save saves clusters information in etcd
func (c Cluster) Save(ctx *config.Context) error {
	val, err := json.Marshal(c)
	if err != nil {
		return errors.New("failed to marshal clusters struct")
	}
	etcd.SaveInEtcd(ctx, getClusterKey(c.ClusterID, c.PoolName), string(val))
	return nil
}

//Delete deletes clusters information from etcd
func (c Cluster) Delete(ctx *config.Context)  {
	etcd.DeleteInEtcd(ctx, getClusterKey(c.ClusterID, c.PoolName))
}

//List gets all the clusters in etcd
func List(ctx *config.Context) ([]Cluster, error)  {
	var clusters []Cluster
	var err error
	d := etcd.GetMultipleWithPrefixFromEtcd(ctx, Cluster_Prefix)
	for _, item := range d {
		var cl Cluster
		err = json.Unmarshal([]byte(item), &cl)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, cl)
	}
	return clusters, nil
}

//ClustersInPool gets all the clusters in a specified pool
func ClustersInPool(ctx *config.Context, poolName string) ([]Cluster,error) {
	var clusters []Cluster
	var err error
	d := etcd.GetMultipleWithPrefixFromEtcd(ctx, getClusterPoolKey(poolName))
	for _, item := range d {
		var cl Cluster
		err = json.Unmarshal([]byte(item), &cl)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, cl)
	}
	return clusters, nil
}

//ClusterByID gets a clusters in a pool with specified ID
func ClusterByID(ctx *config.Context, poolName string, clusterid string) (*Cluster, error) {
	var cluster Cluster
	val := etcd.GetExactFromEtcd(ctx, getClusterKey(clusterid, poolName))
	err := json.Unmarshal([]byte(val), &cluster)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}