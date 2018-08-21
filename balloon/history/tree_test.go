package history

import (
	"fmt"
	"os"
	"testing"

	"github.com/iknite/molletree/encoding/encuint64"
	"github.com/iknite/molletree/hashing"
	"github.com/iknite/molletree/storage"
)

func deleteFile(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Unable to remove db file %s", err)
	}
}

func openBadgerStorage() (storage.Storer, func()) {
	store := storage.NewBadgerStore("/var/tmp/history_store_test.db")
	return store, func() {
		fmt.Println("Cleaning...")
		store.Close()
		deleteFile("/var/tmp/history_store_test.db")
	}
}

func BenchmarkMemoryAdd(b *testing.B) {
	s := storage.NewMemoryStore()

	t := &Tree{version: 0, hasher: hashing.NewSha256(), store: s}
	b.N = 100000

	for i := 0; i < b.N; i++ {
		t.Add(encuint64.ToBytes(uint64(i)))
	}
}

func BenchmarkBadgerAdd(b *testing.B) {
	s, closeF := openBadgerStorage()
	defer closeF()

	t := &Tree{version: 0, hasher: hashing.NewSha256(), store: s}
	b.N = 100000

	for i := 0; i < b.N; i++ {
		t.Add(encuint64.ToBytes(uint64(i)))
	}
}
