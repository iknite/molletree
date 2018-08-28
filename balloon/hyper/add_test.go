package hyper

import (
	"fmt"
	"os"
	"testing"

	"github.com/iknite/molletree/hashing"
	"github.com/iknite/molletree/storage"
	"github.com/stretchr/testify/require"
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

func TestAdd(t *testing.T) {

	testCases := []struct {
		eventDigest      []byte
		expectedRootHash []byte
	}{
		{[]byte{0x0}, []byte{0x0}},
		{[]byte{0x1}, []byte{0x1}},
		{[]byte{0x2}, []byte{0x3}},
		{[]byte{0x3}, []byte{0x0}},
		{[]byte{0x4}, []byte{0x4}},
		{[]byte{0x5}, []byte{0x1}},
		{[]byte{0x6}, []byte{0x7}},
		{[]byte{0x7}, []byte{0x0}},
		{[]byte{0x8}, []byte{0x8}},
		{[]byte{0x9}, []byte{0x1}},
	}

	s, closeF := openBadgerStorage()
	defer closeF()

	h := hashing.NewXor()
	tree := &Tree{
		hasher:     h,
		store:      s.(storage.Cacher),
		cache:      storage.NewMemoryStore(),
		cacheLevel: 1,
	}

	for i, c := range testCases {
		require.Equalf(t, c.expectedRootHash, tree.Add(c.eventDigest), "Incorrect root hash for index %d", i)
	}
}
