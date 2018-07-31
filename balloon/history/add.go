package history

import (
	"github.com/iknite/molletree/encoding/encstring"
)

func (t *Tree) Add(event string) (commitment []byte, digest []byte) {
	// Add a leaf node
	node := &Node{index: t.version, layer: 0, tree: t}
	digest = t.hasher.Cipher(node.Id(), (encstring.ToBytes(event)))
	t.store.Set(node.String(), digest)
	commitment = node.Commitment()

	t.version += 1

	return
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
