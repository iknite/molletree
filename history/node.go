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
	return fmt.Sprintf("%d:%d", n.index, n.layer)
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
	position := n.layer - 1

	return &Node{
		index: n.index + uint64(math.Exp2(float64(position))),
		layer: position,
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

	id := n.String()

	hash, ok := storage.Get(id)
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
		storage.Set(id, hash)
	} else {
		// If bottom nodes are empty warn the upper about it
		tainted = true
	}

	return
}
