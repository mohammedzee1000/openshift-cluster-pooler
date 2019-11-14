package pools

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/mohammedzee1000/openshift-cluster-pooler/pkg/duration"
	"github.com/mohammedzee1000/openshift-cluster-pooler/pkg/utils"

	"github.com/pkg/errors"
)

const (
	//poolExpectedPrefix is the key prefix for pool expected entry
	poolExpectedPrefix = "PoolExpected"
	//poolActualPrefix is the key prefix for pool actual entry
	poolActualPrefix = "PoolActual"
)

//Pool is a pool
type Pool struct {
	Name                        string            `json:"Name"`
	Description                 string            `json:"Description"`
	Size                        int               `json:"Size"`
	MaxSize                     int               `json:"MaxSize"`
	UnusedClusterTimeout        duration.Duration `json:"UnusedClusterTimeout"`
	UsedClusterTimeout          duration.Duration `json:"UsedClusterTimeout"`
	ParallelProvisioning        int               `json:"ParallelProvisioning"`
	ParallelDeProvisioning      int               `json:"ParallelDeProvisioning"`
	ProvisionCommand            string            `json:"ProvisionCommand"`
	DeProvisionCommand          string            `json:"DeProvisionCommand"`
	ForceDeprovisionCommand     string            `json:"ForceDeProvisionCommand"`
	ClusterURLCommand           string            `json:"ClusterUrlCommand"`
	ClusterAdminUserCommand     string            `json:"ClusterAdminUserCommand"`
	ClusterAdminPasswordCommand string            `json:"ClusterAdminPasswordCommand"`
	ClusterCAFilePath           string            `json:"ClusterCAFilePath"`
	ClusterCertFilePath         string            `json:"ClusterCertFilePath"`
	ClusterKeyFilePath          string            `json:"ClusterKeyFilePath"`
	ClusterExtraInfoCommand     string            `json:"ClusterExtraInfoCommand"`
}

//NewEmptyPool creates a new empty pool object
func NewEmptyPool() *Pool {
	return &Pool{}
}

//PoolFromConfig reads a pool from yaml config file
func PoolFromConfig(filename string) (*Pool, error) {
	p := NewEmptyPool()
	err := utils.ReadYamlFile(p, filename)
	if err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal config")
	}
	return p, nil
}

//PoolNameExpectedKey the the poolname expected key
func PoolNameExpectedKey(poolName string) string {
	return fmt.Sprintf("%s-%s", poolExpectedPrefix, poolName)
}

//PoolNameActualKey the the pool_name expected key
func PoolNameActualKey(poolName string) string {
	return fmt.Sprintf("%s-%s", poolActualPrefix, poolName)
}

func PoolByNameWithTransaction(tx *badger.Txn, poolname string, expected bool) (*Pool, error) {
	p := NewEmptyPool()
	var searchkey string
	if !expected {
		searchkey = PoolNameActualKey(poolname)
	} else {
		searchkey = PoolNameExpectedKey(poolname)
	}
	item, err := tx.Get([]byte(searchkey))
	if err != nil {
		return nil, err
	}
	err = item.Value(func(val []byte) error {
		err = json.Unmarshal(val, p)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal pool")
		}
		return nil
	})
	return p, nil
}
