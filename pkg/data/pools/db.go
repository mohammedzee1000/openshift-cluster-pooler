package pools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/database"
)

func GetPoolKey(name string) string {
	return fmt.Sprintf("%s-%s", Pool_Key, name)
}

func GetCleanupPoolKey(name string) string {
	return fmt.Sprintf("%s-%s", Cleanup_Pool_Key, name)
}

//SaveInDB Saves the pool into database
func (p Pool) SaveInDB(ctx *generic.Context, removal bool) error  {
	var key string
	val, err := json.Marshal(p)
	if err != nil {
		return errors.New("failed to marshal pool struct")
	}
	if !removal {
		key = GetPoolKey(p.Name)
	} else {
		key = GetCleanupPoolKey(p.Name)
	}
	database.SaveinKVDB(ctx, key, string(val))
	return nil
}

//deleteInDB deletes the pool from database
func (p Pool) deleteInDB(ctx *generic.Context) {
	database.DeleteInKVDB(ctx, GetPoolKey(p.Name))
}

//Mark for removal marks a pool for garbage collection
func (p Pool) MarkForRemoval(ctx *generic.Context) error {
	val, err := json.Marshal(p)
	if err != nil {
		return errors.New("failed to marshal clusters struct")
	}
	database.SaveinKVDB(ctx, GetCleanupPoolKey(p.Name), string(val))
	p.deleteInDB(ctx)
	return nil
}

func (p Pool) IsMarkedForRemoval(ctx *generic.Context) bool {
	poolNameCleanup := GetCleanupPoolKey(p.Name)
	if database.KeyExistsInKVDB(ctx, poolNameCleanup) {
		return true
	}
	return false
}

//List gets all pools in database
func List(ctx *generic.Context, removal bool) ([]Pool, error)  {
	var pools []Pool
	var key string
	if removal {
		key = Cleanup_Pool_Key
	} else {
		key = Pool_Key
	}
	var err error
	d := database.GetMultipleWithPrefixFromKVDB(ctx, key)
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

//PoolByName gets a pool of specified name. If removal is set, then it gets the pool only
//if it marked for removal
func PoolByName(ctx *generic.Context, name string, removal bool) (*Pool, error)  {
	var p Pool
	var poolkey string
	if removal {
		poolkey = GetCleanupPoolKey(name)
	} else {
		poolkey = GetPoolKey(name)
	}
	val := database.GetExactFromKVDB(ctx, poolkey)
	err := json.Unmarshal([]byte(val), &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}