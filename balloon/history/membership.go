package history

import (
	"github.com/iknite/molletree/encoding/encbytes"
	"github.com/iknite/molletree/hashing"
	"github.com/iknite/molletree/storage"
)

type MembershipProof struct {
	commitment     []byte
	index, version uint64
	store          storage.Storer
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

func (n *Node) AuditPath(version uint64) storage.Storer {

	if n.index > version {
		panic("version is below index")
	}

	store := storage.NewMemoryStore()
	collectAuditPath(store, n.Root(version), n.index, version)

	return store

}

func collectAuditPath(store storage.Storer, node *Node, target, version uint64) {

	if node.layer < 1 {
		// Store the leaf node and end traversing
		store.Set(node.Id(), node.Hash(version))
		return
	}

	right := node.Right()
	left := node.Left()

	if right.index <= target {
		store.Set(left.Id(), left.Hash(version))
		collectAuditPath(store, right, target, version)

	} else {
		if right.index <= version {
			store.Set(right.Id(), right.Hash(version))
		}
		collectAuditPath(store, left, target, version)
	}

}
