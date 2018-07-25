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
	if n.layer < 1 {
		return nil
	}

	return &Node{
		index: n.index,
		layer: n.layer - 1,
		tree:  n.tree,
	}
}

func (n *Node) Right() *Node {
	if n.layer < 1 {
		return nil
	}

	return &Node{
		index: n.index + uint64(math.Exp2(float64(n.layer-1))),
		layer: n.layer - 1,
		tree:  n.tree,
	}
}

func (n *Node) Root() *Node {
	return &Node{
		index: 0,
		layer: treeHeight(n.index),
		tree:  n.tree,
	}
}

func (n *Node) Commitment() []byte {
	rootNode := n.Root()
	hash, _ := rootNode.Hash(n.index)
	return hash
}

func (n *Node) Hash(version uint64) (hash []byte, tainted bool) {

	if n.index > version {
		// if you're trying to get the future, return nil
		return
	}

	nodeString := n.String()

	// REFACTOR: this call is slow if it touches the disk
	hash, ok := n.tree.store.Get(nodeString)
	if ok {
		// you're getting the past so it's cached
		return
	}

	leftN := n.Left()
	leftHash, _ := leftN.Hash(version)

	rightN := n.Right()
	rightHash, rightTainted := rightN.Hash(version)

	hash = n.tree.hasher.Cipher(n.Id(), leftHash, rightHash)

	if rightHash != nil && !rightTainted {
		// is storable when the childrens are complete
		n.tree.store.Set(nodeString, hash)
	} else {
		// If bottom nodes are empty warn the upper about it
		tainted = true
	}

	return
}

func (n *Node) AuditPath() storage.Store {
	rootNode := n.Root()
	store := storage.NewStore()

	collectAuditPath(store, rootNode, n.index)

	hash, _ := n.tree.store.Get(n.String())
	store.Set(n.String(), hash)

	return store
}

func collectAuditPath(store storage.Store, node *Node, version uint64) {

	rightNode := node.Right()
	leftNode := node.Left()

	if rightNode == nil || leftNode == nil {
		return
	}

	if rightNode.index <= version {
		leftHash, _ := leftNode.tree.store.Get(leftNode.String())
		store.Set(leftNode.String(), leftHash)

		collectAuditPath(store, rightNode, version)
	} else {
		collectAuditPath(store, leftNode, version)
	}

}
