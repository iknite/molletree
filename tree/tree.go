package tree

import "math"

type Node struct{}

// Stringer implementation
func (n *Node) String() string {
	return "I'm a node\n"
}

type Tree struct {
	version float64
}

func (t *Tree) height() float64 {
	return math.Ceil(math.Log2(t.version + 1))
}
