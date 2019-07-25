package database

import (
	"github.com/dgraph-io/badger"
	"github.com/mohammedzee1000/openshift-cluster-pool/pkg/generic"
)

// TODO make it retry for database stuff

func HandleError(ctx *generic.Context, err error)  {
		if err != nil {
			ctx.Log.Fatal("Database Ops", err, "failed to connect to data source")
		}
}

//KeyExistsInKVDB checks if paticular key not in DB
func KeyExistsInKVDB(ctx *generic.Context, key string) bool {
	exists := false
	db, err := ctx.NewBadgerConnection()
	if err != nil {
		HandleError(ctx, err)
	}
	defer db.Close()
	err = db.View(func(txn *badger.Txn) error {
		_, err1 := txn.Get([]byte(key))
		if err1 != nil && err1 != badger.ErrKeyNotFound {
			return err1
		}
		exists = true
		return nil
	})
	HandleError(ctx, err)
	return exists
}

//SaveinKVDB saved specified key value pair in database
func SaveinKVDB(ctx *generic.Context, key string, data string) {
	db, err := ctx.NewBadgerConnection()
	if err != nil {
		HandleError(ctx, err)
	}
	defer db.Close()
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(data))
	})
	HandleError(ctx, err)
}

//GetMultipleWithPrefixFromKVDB gets multiple values whose keys match specified prefix in database
func GetMultipleWithPrefixFromKVDB(ctx *generic.Context, keyprefix string) []string {
	var values []string
	db, err := ctx.NewBadgerConnection()
	if err != nil {
		HandleError(ctx, err)
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
	HandleError(ctx, err)
	return values
}

//GetExactFromKVDB gets specific value which matches exact string
func GetExactFromKVDB(ctx *generic.Context, key string) string {
	var value string
	db, err := ctx.NewBadgerConnection()
	HandleError(ctx, err)
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
	HandleError(ctx, err)
	return value
}

//DeleteInKVDB deletes the key specified in database
func DeleteInKVDB(ctx *generic.Context, key string)  {
	db, err := ctx.NewBadgerConnection()
	HandleError(ctx, err)
	defer db.Close()
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	HandleError(ctx, err)
}