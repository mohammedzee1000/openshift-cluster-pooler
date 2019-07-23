package clusters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/config"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/database"
)

func getClusterPoolKey(poolname string) string {
	return fmt.Sprintf("%s-%s", Cluster_Prefix, poolname)
}

func getClusterKey(clusterid string, poolname string) string {
	return fmt.Sprintf("%s-%s-%s", Cluster_Prefix, poolname, clusterid)
}

//Save saves clusters information in database
func (c Cluster) Save(ctx *config.Context) error {
	val, err := json.Marshal(c)
	if err != nil {
		return errors.New("failed to marshal clusters struct")
	}
	database.SaveinKVDB(ctx, getClusterKey(c.ClusterID, c.PoolName), string(val))
	return nil
}

//Delete deletes clusters information from database
func (c Cluster) Delete(ctx *config.Context)  {
	database.DeleteInEtcd(ctx, getClusterKey(c.ClusterID, c.PoolName))
}

//List gets all the clusters in database
func List(ctx *config.Context) ([]Cluster, error)  {
	var clusters []Cluster
	var err error
	d := database.GetMultipleWithPrefixFromKVDB(ctx, Cluster_Prefix)
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
	d := database.GetMultipleWithPrefixFromKVDB(ctx, getClusterPoolKey(poolName))
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
	val := database.GetExactFromKVDB(ctx, getClusterKey(clusterid, poolName))
	err := json.Unmarshal([]byte(val), &cluster)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}