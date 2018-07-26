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
		panic("version is below target")
	}

	store := storage.NewStore()
	collectAuditPath(store, n.Root(version), n.index, version)
	store.Set(n.String(), n.Hash(version))

	return store
}

func collectAuditPath(store storage.Store, node *Node, target, version uint64) {
	if node.layer < 1 {
		return
	}

	rightNode := node.Right()
	leftNode := node.Left()

	if rightNode.index <= target {
		store.Set(leftNode.String(), leftNode.Hash(version))
		collectAuditPath(store, rightNode, target, version)

	} else {
		if rightNode.index <= version {
			store.Set(rightNode.String(), rightNode.Hash(version))
		}
		collectAuditPath(store, leftNode, target, version)
	}
}
