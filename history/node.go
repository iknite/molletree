package history

import (
	"fmt"
	"math"

	"github.com/iknite/bygone-tree/encoding/encuint64"
	"github.com/iknite/bygone-tree/storage"
)

type Node struct {
	index, layer uint64
	tree         *Tree
}

// Stringer implementation
func (n *Node) String() string {
	return fmt.Sprintf("%d|%d", n.index, n.layer)
}

func (n *Node) Id() []byte {
	var b []byte
	b = append(b, encuint64.ToBytes(n.index)...)
	return append(b, encuint64.ToBytes(n.layer)...)
}

func (n *Node) Left() *Node {
	return &Node{
		index: n.index,
		layer: n.layer - 1,
		tree:  n.tree,
	}
}

func (n *Node) Right() *Node {
	return &Node{
		index: n.index + uint64(math.Exp2(float64(n.layer-1))),
		layer: n.layer - 1,
		tree:  n.tree,
	}
}

func (n *Node) Root(version uint64) *Node {
	return &Node{
		index: 0,
		layer: treeHeight(version),
		tree:  n.tree,
	}
}

func (n *Node) Next(indexTarget uint64) *Node {
	right := n.Right()
	if right.index <= indexTarget {
		return right
	} else {
		return n.Left()
	}
}

func (n *Node) Commitment() []byte {
	return n.Root(n.index).Hash(n.index)
}

func (n *Node) Hash(version uint64) []byte {
	hash, _ := n.hash(version)
	return hash
}

func (n *Node) hash(version uint64) (hash []byte, tainted bool) {

	if n.index > version {
		return // if you're trying to get the future, return nil
	}

	key := n.String()

	hash, ok := n.tree.store.Get(key) // TODO: this call is slow if it touches the disk

	if ok {
		return // you're getting the past so it's cached
	}

	rightHash, rightTainted := n.Right().hash(version)

	hash = n.tree.hasher.Cipher(
		n.Id(),
		n.Left().Hash(version),
		rightHash,
	)

	if rightHash != nil && !rightTainted {
		n.tree.store.Set(key, hash) // is storable when the childrens are complete
	} else {
		tainted = true // If bottom nodes are empty warn the upper about it
	}

	return

}

func (n *Node) AuditPath(version uint64) storage.Store {

	if n.index > version {
		panic("version is below index")
	}

	store := storage.NewStore()
	collectAuditPath(store, n.Root(version), n.index, version)

	return store

}

func collectAuditPath(store storage.Store, node *Node, target, version uint64) {

	if node.layer < 1 {
		// Store the leaf node and end traversing
		store.Set(node.String(), node.Hash(version))
		return
	}

	right := node.Right()
	left := node.Left()

	if right.index <= target {
		store.Set(left.String(), left.Hash(version))
		collectAuditPath(store, right, target, version)

	} else {
		if right.index <= version {
			store.Set(right.String(), right.Hash(version))
		}
		collectAuditPath(store, left, target, version)
	}

}

func (n *Node) IncrementalAuditPath(n2 *Node) storage.Store {
	store := storage.NewStore()
	collectIncrementalAuditPath(store, n.Root(n2.index), n.index, n2.index)

	return store

}

func collectIncrementalAuditPath(store storage.Store, node *Node, t1, t2 uint64) {
	if node.layer < 1 {
		store.Set(node.String(), node.Hash(t2))
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
		store.Set(left.String(), left.Hash(t2))
		collectIncrementalAuditPath(store, right, t1, t2)

	} else {
		if right.index <= t2 {
			store.Set(right.String(), right.Hash(t2))
		}
		collectIncrementalAuditPath(store, left, t1, t2)
	}

}
