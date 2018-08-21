// Package bagder implements the store engine interface for
// github.com/dgraph-io/badger/options
package storage

import (
	"bytes"
	"log"

	b "github.com/dgraph-io/badger"
	bo "github.com/dgraph-io/badger/options"
)

type BadgerStore struct {
	db *b.DB
}

func (s BadgerStore) Set(id []byte, data []byte) {
	s.db.Update(func(txn *b.Txn) error {
		return txn.Set(id, data)
	})
}

func (s BadgerStore) Get(id []byte) ([]byte, bool) {
	var value []byte

	err := s.db.View(func(txn *b.Txn) error {
		item, err := txn.Get(id)
		if err != nil {
			return err
		}
		value, _ = item.ValueCopy(value)
		if err != nil {
			return err
		}
		return nil
	})

	switch err {
	case nil:
		return value, true

	case b.ErrKeyNotFound:
		return make([]byte, 0), false

	default:
		return nil, false
	}

}

func (s BadgerStore) GetRange(start, end []byte) [][]byte {
	var leaves [][]byte

	s.db.View(func(txn *b.Txn) error {
		opts := b.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(start); it.Valid(); it.Next() {
			item := it.Item()
			var k []byte
			k = item.KeyCopy(k)
			if bytes.Compare(k, end) > 0 {
				break
			}
			leaves = append(leaves, k)
		}
		return nil
	})

	return leaves
}

func (s BadgerStore) Delete(id []byte) error {
	return s.db.Update(func(txn *b.Txn) error {
		return txn.Delete(id)
	})
}

func (s BadgerStore) Close() error {
	return s.db.Close()
}

func NewBadgerStore(path string) *BadgerStore {
	opts := b.DefaultOptions
	opts.TableLoadingMode = bo.MemoryMap
	opts.Dir = path
	opts.ValueDir = path
	opts.SyncWrites = false
	db, err := b.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	return &BadgerStore{db}
}

func NewBadgerStoreOpts(opts b.Options) (*BadgerStore, *b.DB) {
	db, err := b.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	return &BadgerStore{db}, db

}
