package history

import (
	"math"

	"github.com/iknite/bygone-tree/encoding/encstring"
	"github.com/iknite/bygone-tree/hashing"
	"github.com/iknite/bygone-tree/storage"
)

type Tree struct {
	version uint64
	hasher  hashing.Hasher
}

func NewTree() *Tree {
	return &Tree{
		version: 0,
		hasher:  &hashing.Sha256Hasher{},
	}
}

func treeHeight(version uint64) uint64 {
	return uint64(math.Ceil(math.Log2(float64(version) + 1)))
}

func (t *Tree) Add(event string) []byte {
	// Add a leaf node
	node := &Node{index: t.version, layer: 0, tree: t}
	storage.Set(node.String(), t.hasher.Do((encstring.ToBytes(event))))

	commitment := node.Commitment()

	t.version += 1

	return commitment
}
