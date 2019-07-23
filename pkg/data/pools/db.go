package pools

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/config"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/database"
)

func GetPoolKey(name string) string {
	return fmt.Sprintf("%s-%s", Pool_Key, name)
}

//Save Saves the pool into database
func (p Pool) Save(ctx *config.Context) error  {
	val, err := json.Marshal(p)
	if err != nil {
		return errors.New("failed to marshal clusters struct")
	}
	database.SaveinKVDB(ctx, GetPoolKey(p.Name), string(val))
	return nil
}

//Delete deletes the pool from database
func (p Pool) Delete(ctx *config.Context) {
	database.DeleteInKVDB(ctx, GetPoolKey(p.Name))
}

//List gets all pools in database
func List(ctx *config.Context) ([]Pool, error)  {
	var pools []Pool
	var err error
	d := database.GetMultipleWithPrefixFromKVDB(ctx, Pool_Key)
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
	val := database.GetExactFromKVDB(ctx, GetPoolKey(name))
	err := json.Unmarshal([]byte(val), &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}