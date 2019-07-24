package clusters

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/database"
)

func GetClusterPoolKey(poolname string) string {
	return fmt.Sprintf("%s-%s", Cluster_Prefix, poolname)
}

func GetClusterKey(clusterid string, poolname string) string {
	return fmt.Sprintf("%s-%s-%s", Cluster_Prefix, poolname, clusterid)
}

//Save saves clusters information in database
func (c Cluster) Save(ctx *generic.Context) error {
	val, err := json.Marshal(c)
	if err != nil {
		return errors.New("failed to marshal clusters struct")
	}
	database.SaveinKVDB(ctx, GetClusterKey(c.ClusterID, c.PoolName), string(val))
	return nil
}

//Delete deletes clusters information from database
func (c Cluster) Delete(ctx *generic.Context)  {
	database.DeleteInKVDB(ctx, GetClusterKey(c.ClusterID, c.PoolName))
}

//List gets all the clusters in database
func List(ctx *generic.Context) ([]Cluster, error)  {
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
func ClustersInPool(ctx *generic.Context, poolName string) ([]Cluster,error) {
	var clusters []Cluster
	var err error
	d := database.GetMultipleWithPrefixFromKVDB(ctx, GetClusterPoolKey(poolName))
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
func ClusterByID(ctx *generic.Context, poolName string, clusterid string) (*Cluster, error) {
	var cluster Cluster
	val := database.GetExactFromKVDB(ctx, GetClusterKey(clusterid, poolName))
	err := json.Unmarshal([]byte(val), &cluster)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}
