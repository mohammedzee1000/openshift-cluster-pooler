package pools

import (
	"encoding/json"
	"github.com/dgraph-io/badger"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/clusters"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/database"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
	"github.com/pkg/errors"
	"time"
)

//ActivateCluster activates a cluster, if it is available. This is the only direct db func
//in pool as we need to ensure cluster activation is a transaction
func (p Pool) ActivateCluster(ctx *generic.Context) (*clusters.Cluster, error)  {
	var c clusters.Cluster
	var found bool
	db, err := ctx.NewBadgerConnection()
	if err != nil {
		database.HandleError(ctx, err)
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
				if c.State == clusters.ClusterSuccess {
					if time.Now().Add(3 * time.Minute).Before(p.ExpiresOn(&c)) {
						found = true
						c.State = clusters.ClusterUsed
						c.ActivatedOn = time.Now()
						data, err := json.Marshal(&c)
						if err != nil{
							return err
						}
						return txn.Set([]byte(clusters.GetClusterKey(c.ClusterID, p.Name)), data)
					}
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
		return nil, errors.New("could not get a cluster to activate")
	}
	database.HandleError(ctx, err)
	return &c, nil
}
