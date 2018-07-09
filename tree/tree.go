package tree

import (
	"fmt"
	"math"
)

var storage = make(map[string]string)

type Node struct {
	index, height float64
}

// Stringer implementation
func (n *Node) String() string {
	return fmt.Sprintf("%v:%v", n.index, n.height)
}

type Tree struct {
	length float64
}

func NewTree() *Tree {
	return &Tree{length: 0}
}

func (t *Tree) height() float64 {
	return math.Ceil(math.Log2(t.length + 1))
}

func (t *Tree) Add(event string) *Node {
	// Add a leaf node
	node := &Node{index: t.length, height: 0}
	storage[node.String()] = event

	t.length += 1

	return node
}
