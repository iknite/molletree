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

func (n *Node) Commitment(version uint64) []byte {
	rootNode := &Node{
		index: 0,
		layer: n.tree.Height(),
		tree:  n.tree,
	}

	hash, _ := rootNode.Hash(version)
	return hash
}

func (n *Node) isRootNode() bool {
	return n.index == 0 && n.layer == n.tree.Height()
}

func (n *Node) Hash(version uint64) (hash []byte, tainted bool) {

	// if you're trying to get the future, return nil
	if n.index > version {
		return
	}

	id := n.String()

	hash, ok := storage.Get(id)
	if ok {
		return
	}

	leftN := n.Left()
	leftHash, _ := leftN.Hash(version)

	rightN := n.Right()
	rightHash, rightTainted := rightN.Hash(version)

	hash = n.tree.hasher.Cipher(n.Id(), leftHash, rightHash)

	// is storable when the childrens are complete and you're not a rootNode
	if rightHash != nil && !rightTainted && !n.isRootNode() {
		storage.Set(id, hash)
	} else {
		tainted = true
	}

	return
}
