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

type MembershipProof struct {
	commitment     []byte
	index, version uint64
	store          storage.Store
	hasher         hashing.Hasher
}

func (t *Tree) ProveMembership(commitment []byte, index, version uint64) *MembershipProof {
	store := Node{index: index, layer: 0, tree: t}.AuditPath(version)

	return &MembershipProof{
		commitment: commitment,
		index:      index,
		version:    version,
		store:      store,
		hasher:     t.hasher,
	}
}

func (m *MembershipProof) Verify() bool {
	node := &Node{
		index: m.index,
		layer: 0,
		tree: &Tree{
			version: m.version,
			hasher:  m.hasher,
			store:   m.store,
		},
	}

	commitment := node.Root(version).Hash(version)

	return encbytes.ToString(commitment) == encbytes.ToString(m.commitment)
}

type IncrementalProof struct {
	commitmentA, commitmentB []byte
	indexA, indexB, version  uint64
	storeA, storeB           storage.Store
	hasher                   hashing.Hasher
}

func (t *Tree) ProveIncremental(
	commitmentA []byte, indexA uint64,
	commitmentB []byte, indexB uint64,
	version uint64,

) *IncrementalProof {

	storeA := Node{index: indexA, layer: 0, tree: t}.AuditPath(version)
	storeB := Node{index: indexB, layer: 0, tree: t}.AuditPath(version)

	return &IncrementalProof{
		commitmentA, commitmentB,
		indexA, indexB, version,
		storeA, storeB,
		t.hasher,
	}
}

func (i *IncrementalProof) Verify() bool {
	rootNodeA := &Node{
		index: i.indexA,
		layer: 0,
		tree: &Tree{
			version: i.version,
			hasher:  i.hasher,
			store:   i.storeA,
		},
	}.Root(i.indexA)

	i.storeB.Set(rootNodeA.String(), rootNodeA.Hash(version))

	node := &Node{
		index: i.indexB,
		layer: 0,
		tree: &Tree{
			version: i.version,
			hasher:  i.hasher,
			store:   i.storeB,
		},
	}

	commitment := node.Root(version).Hash(version)

	return encbytes.ToString(commitment) == encbytes.ToString(m.commitment)
}
