package hyper

import (
	"fmt"

	"github.com/iknite/molletree/encoding/encuint64"
	"github.com/iknite/molletree/storage"
)

type Node struct {
	base    []byte
	layer   uint64
	tree    *Tree
	opStore storage.Merger
}

// Stringer implementation
func (n *Node) String() string {
	return fmt.Sprintf("%08b|%d", n.base, n.layer)
}

func (n *Node) Id() []byte {
	var b []byte
	b = append(b, n.base...)
	return append(b, encuint64.ToBytes(n.layer)...)
}

func (n *Node) Height() uint64 {
	return n.tree.hasher.Len()
}

func (n *Node) Left() *Node {
	var b []byte

	return &Node{
		base:    append(b, n.base...),
		layer:   n.layer - 1,
		tree:    n.tree,
		opStore: n.opStore,
	}
}

func (n *Node) Right() *Node {
	var b []byte
	b = append(b, n.base...)
	bitSet(b, n.Height()-n.layer)

	return &Node{
		base:    b,
		layer:   n.layer - 1,
		tree:    n.tree,
		opStore: n.opStore,
	}
}

func (n *Node) Root() *Node {
	height := n.Height()
	return &Node{
		base:    make([]byte, (height)/8),
		layer:   height,
		tree:    n.tree,
		opStore: n.opStore,
	}
}

func (n *Node) Next(target []byte) *Node {
	if bitIsSet(target, (n.Height() - max(1, n.layer))) {
		return n.Right()
	} else {
		return n.Left()
	}
}

func (n *Node) Store() (s storage.Storer) {
	if n.layer >= n.tree.cacheLevel {
		s = n.tree.cache
	} else {
		s = n.opStore
	}
	return
}

func bitIsSet(bits []byte, i uint64) bool { return bits[i/8]&(1<<uint(7-i%8)) != 0 }
func bitSet(bits []byte, i uint64)        { bits[i/8] |= 1 << uint(7-i%8) }

func max(a, b uint64) uint64 {
	if a > b {
		return a
	} else {
		return b
	}
}
