package database

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/clusters"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/data/pools"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"github.com/pkg/errors"
)
//ActivateClusterInPool a cluster, if it is available. This is the only direct db func
//in pool as we need to ensure cluster activation is a transaction
//in here because its only pool level op that needs full fledged lock
func ActivaeClusterInPool(ctx *generic.Context, p *pools.Pool) (*clusters.Cluster, error)  {
	var c clusters.Cluster
	var found bool
	db, err := ctx.NewBadgerConnection()
	if err != nil {
		HandleError(ctx, err)
	}
	defer db.Close()
	// cluster activation transaction
	err = db.Update(func(txn *badger.Txn) error {
		//collect all clusters in pool
		keyprefix := clusters.GetClusterPoolKey(p.Name)
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(keyprefix)
		//iterate over clusters in db, trying to find one which is in success state and can be activated
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			// work with the current item
			err := item.Value(func(v []byte) error {
				err = json.Unmarshal(v, &c)
				if err != nil {
					return err
				}
				// if one of the clusters we are iterating can be activated, activate it and found is true
				if c.State == clusters.State_Success {
					found = true
					c.State = clusters.State_Used
					data, err := json.Marshal(&c)
					if err != nil{
						return err
					}
					return txn.Set([]byte(clusters.GetClusterKey(c.ClusterID, p.Name)), data)
				}
				return nil
			})
			if err != nil {
				return err
			}
			//if we have found what we want, return nil breaking out of the
			if found {
				return nil
			}
		}
		return nil
	})
	// if not found return an error
	if !found {
		return nil, errors.New("could not get a cluster to activae")
	}
	HandleError(ctx, err)
	return &c, nil
}
