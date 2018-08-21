package history

import (
	"fmt"

	"github.com/iknite/molletree/encoding/encuint64"
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

func (n *Node) Capacity() uint64 {
	return n.index + 1<<n.layer - 1
}

func (n *Node) Right() *Node {
	return &Node{
		index: n.index + (1 << (n.layer - 1)),
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
