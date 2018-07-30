package history

import (
	"math"

	"github.com/iknite/molletree/encoding/encbytes"
	"github.com/iknite/molletree/encoding/encstring"
	"github.com/iknite/molletree/hashing"
	"github.com/iknite/molletree/storage"
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
	node := &Node{index: index, layer: 0, tree: t}
	store := node.AuditPath(version)

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

	commitment := node.Root(m.version).Hash(m.version)

	return encbytes.ToString(commitment) == encbytes.ToString(m.commitment)
}

type IncrementalProof struct {
	start, end uint64
	store      storage.Store
	hasher     hashing.Hasher
}

func (t *Tree) ProveIncremental(start, end uint64) *IncrementalProof {
	startNode := &Node{index: start, layer: 0, tree: t}
	endNode := &Node{index: end, layer: 0, tree: t}

	return &IncrementalProof{
		start:  start,
		end:    end,
		store:  startNode.IncrementalAuditPath(endNode),
		hasher: t.hasher,
	}
}

func (p *IncrementalProof) Verify(startHash, endHash []byte) bool {
	t := &Tree{version: p.end, hasher: p.hasher, store: p.store}
	startNode := &Node{index: p.start, layer: 0, tree: t}
	endNode := &Node{index: p.end, layer: 0, tree: t}

	startNode.Commitment()
	startCommitment := startNode.Commitment()
	endCommitment := endNode.Commitment()

	return encbytes.ToString(startCommitment) == encbytes.ToString(startHash) &&
		encbytes.ToString(endCommitment) == encbytes.ToString(endHash)
}
