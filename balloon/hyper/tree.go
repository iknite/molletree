package hyper

import (
	"github.com/iknite/molletree/hashing"
	"github.com/iknite/molletree/storage"
)

type Tree struct {
	hasher     hashing.Hasher
	store      storage.Cacher
	cache      storage.Storer
	cacheLevel uint64
}

func NewTree() *Tree {
	h := hashing.NewSha256()

	return &Tree{
		hasher:     h,
		store:      storage.NewBadgerStore("/var/tmp/balloon-hyper"),
		cache:      storage.NewMemoryStore(),
		cacheLevel: h.Len() / 2,
	}
}
