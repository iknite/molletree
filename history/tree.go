package history

import (
	"math"

	"github.com/iknite/bygone-tree/encoding/encbytes"
	"github.com/iknite/bygone-tree/encoding/encstring"
	"github.com/iknite/bygone-tree/hashing"
	"github.com/iknite/bygone-tree/storage"
)

type Tree struct {
	version uint64
	hasher  hashing.Hasher
	store   storage.Store
}

func NewTree() *Tree {
	return &Tree{
		version: 0,
		hasher:  &hashing.Sha256Hasher{},
		store:   storage.NewStore(),
	}
}

func treeHeight(version uint64) uint64 {
	return uint64(math.Ceil(math.Log2(float64(version) + 1)))
}

func (t *Tree) Add(event string) []byte {
	// Add a leaf node
	node := &Node{index: t.version, layer: 0, tree: t}
	t.store.Set(node.String(), t.hasher.Do((encstring.ToBytes(event))))
	commitment := node.Commitment()

	t.version += 1

	return commitment
}

type Proof struct {
	commitment []byte
	version    uint64
	store      storage.Store
	hasher     hashing.Hasher
}

func (t *Tree) MembershipProof(commitment []byte, index uint64, version uint64) *Proof {
	node := &Node{index: index, layer: 0, tree: t}
	audithpath := node.AuditPath(version)

	return &Proof{
		commitment: commitment,
		version:    version,
		store:      audithpath,
		hasher:     t.hasher,
	}
}

func (p *Proof) Verify() bool {
	t := &Tree{version: p.version, hasher: p.hasher, store: p.store}
	node := &Node{index: t.version, layer: 0, tree: t}

	return encbytes.ToString(node.Commitment()) == encbytes.ToString(p.commitment)
}
