package history

import (
	"github.com/iknite/molletree/encoding/encbytes"
	"github.com/iknite/molletree/hashing"
	"github.com/iknite/molletree/storage"
)

type IncrementalProof struct {
	start, end uint64
	store      storage.Storer
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

func (n *Node) IncrementalAuditPath(n2 *Node) storage.Storer {
	store := storage.NewMemoryStore()
	collectIncrementalAuditPath(store, n.Root(n2.index), n.index, n2.index)

	return store

}

func collectIncrementalAuditPath(store storage.Storer, node *Node, t1, t2 uint64) {
	if node.layer < 1 {
		store.Set(node.Id(), node.Hash(t2))
		return
	}

	targetN1 := node.Next(t1)
	targetN2 := node.Next(t2)

	if targetN1.index != targetN2.index {
		// if there is a split between paths, compute the paths independtly
		collectAuditPath(store, targetN1, t1, t2)
		collectAuditPath(store, targetN2, t2, t2)
		return
	}

	right := node.Right()
	left := node.Left()

	if right.index == targetN1.index {
		store.Set(left.Id(), left.Hash(t2))
		collectIncrementalAuditPath(store, right, t1, t2)

	} else {
		if right.index <= t2 {
			store.Set(right.Id(), right.Hash(t2))
		}
		collectIncrementalAuditPath(store, left, t1, t2)
	}

}
