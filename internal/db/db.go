package db

import (
	"github.com/dgraph-io/badger"
)

type DB struct {
	badgerDB *badger.DB
}

var instance *DB

func Open(directory string) (*DB, error) {
	db, err := badger.Open(badger.DefaultOptions(directory))

	if err == nil {
		instance = &DB{
			badgerDB: db,
		}
		return instance, nil
	} else {
		return nil, err
	}
}

func GetKvDb() *badger.DB {
	return instance.badgerDB
}

func (d *DB) Close() {
	d.badgerDB.Close()
}
