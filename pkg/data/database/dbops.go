package database

import (
	"github.com/dgraph-io/badger"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/config"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/logging"
)

// TODO make it retry for database stuff

func HandleError(err error)  {
		if err != nil {
			logging.Fatal("Database Ops", "failed to connect to data source : %s", err.Error())
		}
}

//SaveinKVDB saved specified key value pair in database
func SaveinKVDB(ctx *config.Context, key string, data string) {
	db, err := ctx.NewBadgerConnection()
	if err != nil {
		HandleError(err)
	}
	defer db.Close()
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(data))
	})
	HandleError(err)
}

//GetMultipleWithPrefixFromKVDB gets multiple values whose keys match specified prefix in database
func GetMultipleWithPrefixFromKVDB(ctx *config.Context, keyprefix string) []string {
	var values []string
	db, err := ctx.NewBadgerConnection()
	if err != nil {
		HandleError(err)
	}
	defer db.Close()
	err = db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(keyprefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				values = append(values, string(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	HandleError(err)
	return values
}

//GetExactFromKVDB gets specific value which matches exact string
func GetExactFromKVDB(ctx *config.Context, key string) string {
	var value string
	db, err := ctx.NewBadgerConnection()
	HandleError(err)
	defer db.Close()
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
	})
	HandleError(err)
	return value
}

//DeleteInKVDB deletes the key specified in database
func DeleteInKVDB(ctx *config.Context, key string)  {
	db, err := ctx.NewBadgerConnection()
	HandleError(err)
	defer db.Close()
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	HandleError(err)
}