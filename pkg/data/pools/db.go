package pools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/config"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/etcd"
)

func getPoolKey(name string) string {
	return fmt.Sprintf("%s-%s", Pool_Key, name)
}

//Save Saves the pool into etcd
func (p Pool) Save(ctx *config.Context) error  {
	val, err := json.Marshal(p)
	if err != nil {
		return errors.New("failed to marshal clusters struct")
	}
	etcd.SaveInEtcd(ctx, getPoolKey(p.Name), string(val))
	return nil
}

//Delete deletes the pool from etcd
func (p Pool) Delete(ctx *config.Context) {
	etcd.DeleteInEtcd(ctx, getPoolKey(p.Name))
}

//List gets all pools in etcd
func List(ctx *config.Context) ([]Pool, error)  {
	var pools []Pool
	var err error
	d := etcd.GetMultipleWithPrefixFromEtcd(ctx, Pool_Key)
	for _, item := range d {
		var p Pool
		err = json.Unmarshal([]byte(item), &p)
		if err != nil {
			return nil, err
		}
		pools = append(pools, p)
	}
	return pools, nil
}

//PoolByName gets a pool of specified name
func PoolByName(ctx *config.Context, name string) (*Pool, error)  {
	var p Pool
	val := etcd.GetExactFromEtcd(ctx, getPoolKey(name))
	err := json.Unmarshal([]byte(val), &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}