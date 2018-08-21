package history

import (
	"math"

	"github.com/iknite/molletree/hashing"
	"github.com/iknite/molletree/storage"
)

type Tree struct {
	version uint64
	hasher  hashing.Hasher
	store   storage.Storer
}

func NewTree() *Tree {
	return &Tree{
		version: 0,
		hasher:  hashing.NewSha256(),
		store:   storage.NewMemoryStore(),
	}
}

func treeHeight(version uint64) uint64 {
	return uint64(math.Ceil(math.Log2(float64(version) + 1)))
}
